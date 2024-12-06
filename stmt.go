package sqlx

import (
	"context"
	"database/sql"
	"reflect"
	"unsafe"
)

type IArgs interface {
	Args() []any
}

type IScanableStruct interface {
	Fields(stmttype reflect.Type) []ISqlField
}

type commonStmt[A IArgs, S any] struct {
	Sql  string
	stmt *sql.Stmt
}

func (cs *commonStmt[A, S]) self() *S {
	return (*S)(unsafe.Pointer(cs))
}

func (cs *commonStmt[A, S]) Prepare(ctx context.Context, db *sql.DB) *S {
	sv, err := db.PrepareContext(ctx, cs.Sql)
	if err != nil {
		panic(err)
	}
	cs.stmt = sv
	return cs.self()
}

type queryStmt[A IArgs, R IScanableStruct] struct {
	commonStmt[A, queryStmt[A, R]]
	constructor func(ctx context.Context) R
	maxrows     int
	selftype    reflect.Type
}

func QueryStmt[A IArgs, R IScanableStruct](sql string) *queryStmt[A, R] {
	return &queryStmt[A, R]{
		commonStmt: commonStmt[A, queryStmt[A, R]]{Sql: sql},
		selftype:   reflect.TypeOf((*queryStmt[A, R])(nil)),
	}
}

func (query *queryStmt[A, R]) Constructor(fnc func(ctx context.Context) R) *queryStmt[A, R] {
	query.constructor = fnc
	return query
}

func (query *queryStmt[A, R]) MaxRows(maxrows int) *queryStmt[A, R] {
	query.maxrows = maxrows
	return query
}

func (query *queryStmt[A, R]) Query(ctx context.Context, args A) ([]R, error) {
	rows, err := query.stmt.QueryContext(ctx, args.Args()...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var vs []R
	var ptrs []any
	for rows.Next() {
		ele := query.constructor(ctx)
		fields := ele.Fields(query.selftype)

		ptrs = ptrs[:0]
		for _, f := range fields {
			ptrs = append(ptrs, f.ScanPtr())
		}

		err := rows.Scan(ptrs...)
		if err != nil {
			return nil, err
		}
		vs = append(vs, ele)
	}
	return vs, nil
}

type execStmt[A IArgs] struct {
	commonStmt[A, execStmt[A]]
}

func (exec *execStmt[A]) Exec(ctx context.Context, args A) (sql.Result, error) {
	return exec.stmt.ExecContext(ctx, args.Args()...)
}

func ExecStmt[A IArgs](sql string) *execStmt[A] {
	return &execStmt[A]{
		commonStmt[A, execStmt[A]]{Sql: sql},
	}
}
