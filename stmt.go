package sqlx

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"sync"
	"unsafe"
)

type IScanable interface {
	FieldPtrs() []any
}

type IArgs interface {
	NamedArgs() []sql.NamedArg
}

var (
	typeofIScanable = reflect.TypeOf((*IScanable)(nil)).Elem()
	typeofIArgs     = reflect.TypeOf((*IArgs)(nil)).Elem()
)

type _CommonStmt[Args any, Self any] struct {
	sql   string
	stmts map[*sql.DB]*sql.Stmt
	_stmt *sql.Stmt

	isIArgs bool

	argvGetters []func(ptr unsafe.Pointer) any
}

func (cs *_CommonStmt[Args, Self]) self() *Self {
	return (*Self)(unsafe.Pointer(cs))
}

func (cs *_CommonStmt[Args, Self]) init(sqltxt string) {
	ti := gettypeinfo[Args](nil)
	if len(ti.fields) < 1 {
		panic(fmt.Errorf("sqlx: empty fields on type, %s", ti.modeltype))
	}
	cs.sql = sqltxt
	if reflect.PointerTo(ti.modeltype).Implements(typeofIArgs) {
		cs.isIArgs = true
	}

	for idx := range ti.fields {
		fmp := &ti.fields[idx]
		cs.argvGetters = append(cs.argvGetters, func(ptr unsafe.Pointer) any {
			fptrv := reflect.NewAt(fmp.Field.Type, unsafe.Add(ptr, fmp.Offset))
			return sql.Named(fmp.Metainfo.Name, fptrv.Elem().Interface())
		})
	}
}

func (cs *_CommonStmt[Args, Self]) Prepare(ctx context.Context, dbs ...*sql.DB) error {
	if len(dbs) < 1 {
		panic(fmt.Errorf("sqlx: empty databases"))
	}

	if len(dbs) == 1 {
		sv, err := dbs[0].PrepareContext(ctx, cs.sql)
		if err != nil {
			return err
		}
		cs._stmt = sv
		return nil
	}

	cs.stmts = make(map[*sql.DB]*sql.Stmt)
	for _, db := range dbs {
		sv, err := db.PrepareContext(ctx, cs.sql)
		if err != nil {
			return err
		}
		cs.stmts[db] = sv
	}
	return nil
}

func (cs *_CommonStmt[Args, Self]) MustPrepare(ctx context.Context, dbs ...*sql.DB) *Self {
	if err := cs.Prepare(ctx, dbs...); err != nil {
		panic(err)
	}
	return cs.self()
}

func (cs *_CommonStmt[Args, Self]) getsv(ctx context.Context) *sql.Stmt {
	if cs._stmt != nil {
		return cs._stmt
	}
	dbv := ctx.Value(ctxKeyForDb)
	if dbv == nil {
		panic(ErrNoDB)
	}
	return cs.stmts[dbv.(*sql.DB)]
}

func (cs *_CommonStmt[Args, Self]) expandArgs(v *Args) []any {
	var qargs []any
	if cs.isIArgs {
		nargs := ((any)(v).(IArgs)).NamedArgs()
		for _, v := range nargs {
			qargs = append(qargs, v)
		}
	} else {
		ptr := unsafe.Pointer(v)
		for _, getter := range cs.argvGetters {
			qargs = append(qargs, getter(ptr))
		}
	}
	return qargs
}

type _SelectStmt[Args any, Scanable any] struct {
	_CommonStmt[Args, _SelectStmt[Args, Scanable]]
	scanfields  []_Field
	constructor func() *Scanable
	lengthhint  int

	isIScanable bool
	lock        sync.RWMutex
	fptrGetters []func(ins unsafe.Pointer) any
}

func SelectStmt[Args any, Scanable any](sql string, constructor func() *Scanable) *_SelectStmt[Args, Scanable] {
	var ti = gettypeinfo[Scanable](nil)
	if len(ti.fields) < 1 {
		panic(fmt.Errorf("sqlx: empty fields on type, %s", ti.modeltype))
	}

	obj := &_SelectStmt[Args, Scanable]{
		constructor: constructor,
		scanfields:  ti.fields,
		isIScanable: reflect.PointerTo(ti.modeltype).Implements(typeofIScanable),
	}
	obj.init(sql)
	return obj
}

func (stmt *_SelectStmt[Args, Scanable]) QueryLengthHint(hint int) *_SelectStmt[Args, Scanable] {
	stmt.lengthhint = hint
	return stmt
}

func (stmt *_SelectStmt[Args, Scanable]) mkPtrGetters(rows *sql.Rows) error {
	stmt.lock.Lock()
	defer stmt.lock.Unlock()

	names, err := rows.Columns()
	if err != nil {
		return err
	}

	var fs []*_Field
	for _, name := range names {
		found := false
		for idx := range stmt.scanfields {
			fp := &stmt.scanfields[idx]
			if fp.Metainfo.Name == name {
				fs = append(fs, fp)
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("sqlx: can not find")
		}
	}

	for _, fp := range fs {
		stmt.fptrGetters = append(stmt.fptrGetters, func(ins unsafe.Pointer) any {
			return reflect.NewAt(fp.Field.Type, unsafe.Add(ins, fp.Offset)).Interface()
		})
	}
	return nil
}

func (stmt *_SelectStmt[Args, Scanable]) ensurePtrGetters(rows *sql.Rows) error {
	stmt.lock.RLock()
	if stmt.fptrGetters != nil {
		stmt.lock.RUnlock()
		return nil
	}
	stmt.lock.RUnlock()
	return stmt.mkPtrGetters(rows)
}

var ErrNoDB = errors.New("sqlx: can not get *sql.DB from context")

func (stmt *_SelectStmt[Args, Scanable]) scanByInterface(rows *sql.Rows) ([]*Scanable, error) {
	var vec []*Scanable
	if stmt.lengthhint > 0 {
		vec = make([]*Scanable, 0, stmt.lengthhint)
	}
	for rows.Next() {
		var eleptr = stmt.constructor()
		var ifele = ((any)(eleptr)).(IScanable)

		var ptrs = ifele.FieldPtrs()
		err := rows.Scan(ptrs...)
		if err != nil {
			return nil, err
		}
		vec = append(vec, eleptr)
	}
	return vec, nil
}

func (stmt *_SelectStmt[Args, Scanable]) QueryMany(ctx context.Context, args *Args) ([]*Scanable, error) {
	var sv = stmt.getsv(ctx)

	txv := ctx.Value(ctxKeyForTx)
	if txv != nil {
		tx := txv.(*sql.Tx)
		sv = tx.StmtContext(ctx, sv)
	}

	rows, err := sv.QueryContext(ctx, stmt.expandArgs(args)...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if stmt.isIScanable {
		return stmt.scanByInterface(rows)
	}

	if err = stmt.ensurePtrGetters(rows); err != nil {
		return nil, err
	}

	var vec []*Scanable
	if stmt.lengthhint > 0 {
		vec = make([]*Scanable, 0, stmt.lengthhint)
	}

	var tmps []any = make([]any, len(stmt.fptrGetters))
	for rows.Next() {
		ele := stmt.constructor()
		ptrnum := unsafe.Pointer(ele)
		for idx, geter := range stmt.fptrGetters {
			tmps[idx] = geter(ptrnum)
		}
		err = rows.Scan(tmps...)
		if err != nil {
			return nil, err
		}
		vec = append(vec, ele)
	}
	return vec, nil
}

func (stmt *_SelectStmt[Args, Scanable]) MustQueryMany(ctx context.Context, args *Args) []*Scanable {
	return must(stmt.QueryMany(ctx, args))
}

func (stmt *_SelectStmt[Args, Scanable]) QueryOne(ctx context.Context, args *Args) (*Scanable, error) {
	vs, err := stmt.QueryMany(ctx, args)
	if err != nil {
		return nil, err
	}
	if len(vs) < 1 {
		return nil, sql.ErrNoRows
	}
	return vs[0], nil
}

func (stmt *_SelectStmt[Args, Scanable]) MustQueryOne(ctx context.Context, args *Args) *Scanable {
	return must(stmt.QueryOne(ctx, args))
}
