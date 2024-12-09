package sqlx

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"sync"
	"unsafe"
)

type selectStmt[S any] struct {
	sql         string
	stmt        *sql.Stmt
	fields      []_Field
	constructor func() *S
	ptrGetters  []func(ins int64) any
	lengthhint  int
	once        sync.Once
}

func SelectStmt[S any](sql string, constructor func() *S) *selectStmt[S] {
	var fs []_Field
	var stype = reflect.TypeOf((*S)(nil)).Elem()

	tv, ok := tables[stype]
	if ok {
		tvv := tv.(*_Table[S])
		fs = tvv.fields
	} else {
		var ptr = Mptr[S]()
		addfield(&fs, ptr, int64(uintptr(unsafe.Pointer(ptr))))
	}
	return &selectStmt[S]{
		sql:         sql,
		fields:      fs,
		constructor: constructor,
	}
}

func (stmt *selectStmt[S]) QueryLengthHint(hint int) *selectStmt[S] {
	stmt.lengthhint = hint
	return stmt
}

func (stmt *selectStmt[S]) mkPtrGetters(rows *sql.Rows) error {
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
		// TODO test this and benchmark reflect.Value.FieldByIndex().Addr().Interface()
		stmt.ptrGetters = append(stmt.ptrGetters, func(ins int64) any {
			var ptr = unsafe.Pointer(uintptr(ins + fp.Offset))
			return reflect.NewAt(fp.Field.Type, ptr).Interface()
		})
	}
	return nil
}

func (stmt *selectStmt[S]) Query(ctx context.Context, args ...any) ([]*S, error) {
	rows, err := stmt.stmt.QueryContext(ctx, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stmt.once.Do(func() { err = stmt.mkPtrGetters(rows) })
	if err != nil {
		return nil, err
	}

	var vec []*S
	if stmt.lengthhint > 0 {
		vec = make([]*S, 0, stmt.lengthhint)
	}

	var tmps []any = make([]any, len(stmt.ptrGetters))
	for rows.Next() {
		ele := stmt.constructor()
		ptrnum := int64(uintptr(unsafe.Pointer(ele)))
		for idx, geter := range stmt.ptrGetters {
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
