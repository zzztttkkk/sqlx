package sqlx

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"sync"
	"unsafe"

	"github.com/zzztttkkk/lion"
)

type _CommonStmt[Args any, Self any] struct {
	sql   string
	stmts map[*sql.DB]*sql.Stmt
	_stmt *sql.Stmt

	isEmptyArgs bool

	argvGetters []func(ptr unsafe.Pointer) any
}

func (cs *_CommonStmt[Args, Self]) self() *Self {
	return (*Self)(unsafe.Pointer(cs))
}

func (cs *_CommonStmt[Args, Self]) init(sqltxt string) {
	ti := lion.TypeInfoOf[Args, DdlOptions]()
	if len(ti.Fields) < 1 {
		cs.isEmptyArgs = true
	}
	cs.sql = sqltxt

	if !cs.isEmptyArgs {
		for idx := range ti.Fields {
			fmp := &ti.Fields[idx]
			cs.argvGetters = append(cs.argvGetters, func(ptr unsafe.Pointer) any {
				fv := fmp.PtrGetter()(ptr)
				return sql.Named(fmp.Name(), reflect.ValueOf(fv).Elem().Interface())
			})
		}
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

func (cs *_CommonStmt[Args, Self]) getdbsv(ctx context.Context) *sql.Stmt {
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
	if cs.isEmptyArgs {
		return nil
	}

	var qargs = make([]any, 0, len(cs.argvGetters))
	ptr := unsafe.Pointer(v)
	for _, getter := range cs.argvGetters {
		qargs = append(qargs, getter(ptr))
	}
	return qargs
}

type _SelectStmt[Args any, Scanable any] struct {
	_CommonStmt[Args, _SelectStmt[Args, Scanable]]
	scanfields  []lion.Field[DdlOptions]
	constructor func() *Scanable
	lengthhint  int

	lock        sync.RWMutex
	fptrGetters []func(ins unsafe.Pointer) any
}

func SelectStmt[Args any, Scanable any](sql string, constructor func() *Scanable) *_SelectStmt[Args, Scanable] {
	var ti = lion.TypeInfoOf[Scanable, DdlOptions]()
	if len(ti.Fields) < 1 {
		panic(fmt.Errorf("sqlx: empty fields on type %s", ti.GoType))
	}

	obj := &_SelectStmt[Args, Scanable]{
		constructor: constructor,
		scanfields:  ti.Fields,
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

	var fs []*lion.Field[DdlOptions]
	for _, name := range names {
		found := false
		for idx := range stmt.scanfields {
			fp := &stmt.scanfields[idx]
			if fp.Name() == name {
				fs = append(fs, fp)
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("sqlx: can not find dest ptr for `%s`", name)
		}
	}

	for _, fp := range fs {
		stmt.fptrGetters = append(stmt.fptrGetters, fp.PtrGetter())
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

func (stmt *_SelectStmt[Args, Scanable]) QueryMany(ctx context.Context, args *Args) ([]*Scanable, error) {
	var sv = stmt.getdbsv(ctx)

	txv := ctx.Value(ctxKeyForTx)
	if txv != nil {
		sv = (txv.(*sql.Tx)).StmtContext(ctx, sv)
		// we do not need to close this stmt
		// sql.Tx will do close when commit/rollback ==> go:src/database/sql/sql.go#2270 func closePrepared
	}

	rows, err := sv.QueryContext(ctx, stmt.expandArgs(args)...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

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
		eleptr := unsafe.Pointer(ele)
		for idx, getter := range stmt.fptrGetters {
			tmps[idx] = getter(eleptr)
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

type _ExecStmt[Args any] struct {
	_CommonStmt[Args, _ExecStmt[Args]]
}

func ExecStmt[Args any](sql string) *_ExecStmt[Args] {
	obj := &_ExecStmt[Args]{}
	obj.init(sql)
	return obj
}

func (stmt *_ExecStmt[Args]) Exec(ctx context.Context, args *Args) (sql.Result, error) {
	var sv = stmt.getdbsv(ctx)
	txv := ctx.Value(ctxKeyForTx)
	if txv != nil {
		sv = (txv.(*sql.Tx)).StmtContext(ctx, sv)
	}
	return sv.ExecContext(ctx, stmt.expandArgs(args)...)
}
