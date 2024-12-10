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

type commonStmt[S any] struct {
	sql   string
	stmts map[*sql.DB]*sql.Stmt
	_stmt *sql.Stmt
}

func (cs *commonStmt[S]) self() *S {
	return (*S)(unsafe.Pointer(cs))
}

func (cs *commonStmt[S]) Prepare(ctx context.Context, dbs ...*sql.DB) error {
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

func (cs *commonStmt[S]) MustPrepare(ctx context.Context, dbs ...*sql.DB) *S {
	if err := cs.Prepare(ctx, dbs...); err != nil {
		panic(err)
	}
	return cs.self()
}

func (cs *commonStmt[S]) getsv(ctx context.Context) *sql.Stmt {
	if cs._stmt != nil {
		return cs._stmt
	}
	dbv := ctx.Value(ctxKeyForDb)
	if dbv == nil {
		panic(ErrNoDB)
	}
	return cs.stmts[dbv.(*sql.DB)]
}

type selectStmt[S any] struct {
	commonStmt[selectStmt[S]]
	fields      []_Field
	constructor func() *S
	lengthhint  int

	lock        sync.RWMutex
	fptrGetters []func(ins unsafe.Pointer) any
}

func SelectStmt[S any](sql string, constructor func() *S) *selectStmt[S] {
	var ti = gettypeinfo[S](nil)
	if len(ti.fields) < 1 {
		panic(fmt.Errorf("sqlx: empty fields on type, %s", ti.modeltype))
	}

	obj := &selectStmt[S]{
		constructor: constructor,
		fields:      ti.fields,
	}
	obj.sql = sql
	return obj
}

func (stmt *selectStmt[S]) QueryLengthHint(hint int) *selectStmt[S] {
	stmt.lengthhint = hint
	return stmt
}

func (stmt *selectStmt[S]) mkPtrGetters(rows *sql.Rows) error {
	stmt.lock.Lock()
	defer stmt.lock.Unlock()

	names, err := rows.Columns()
	if err != nil {
		return err
	}

	var fs []*_Field
	for _, name := range names {
		found := false
		for idx := range stmt.fields {
			fp := &stmt.fields[idx]
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

func (stmt *selectStmt[S]) ensurePtrGetters(rows *sql.Rows) error {
	stmt.lock.RLock()
	if stmt.fptrGetters != nil {
		stmt.lock.RUnlock()
		return nil
	}
	stmt.lock.RUnlock()
	return stmt.mkPtrGetters(rows)
}

var ErrNoDB = errors.New("sqlx: can not get *sql.DB from context")

func (stmt *selectStmt[S]) QueryMany(ctx context.Context, args ...any) ([]*S, error) {
	var sv = stmt.getsv(ctx)

	txv := ctx.Value(ctxKeyForTx)
	if txv != nil {
		tx := txv.(*sql.Tx)
		sv = tx.StmtContext(ctx, sv)
	}

	rows, err := sv.QueryContext(ctx, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if err = stmt.ensurePtrGetters(rows); err != nil {
		return nil, err
	}

	var vec []*S
	if stmt.lengthhint > 0 {
		vec = make([]*S, 0, stmt.lengthhint)
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

func (stmt *selectStmt[S]) MustQueryMany(ctx context.Context, args ...any) []*S {
	return must(stmt.QueryMany(ctx, args...))
}

func (stmt *selectStmt[S]) QueryOne(ctx context.Context, args ...any) (*S, error) {
	vs, err := stmt.QueryMany(ctx, args...)
	if err != nil {
		return nil, err
	}
	if len(vs) < 1 {
		return nil, sql.ErrNoRows
	}
	return vs[0], nil
}

func (stmt *selectStmt[S]) MustQueryOne(ctx context.Context, args ...any) *S {
	return must(stmt.QueryOne(ctx, args...))
}
