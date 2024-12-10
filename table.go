package sqlx

import (
	"reflect"
)

type _Table[T any] struct {
	*_TypeInfo[T]

	ddl     *_Ddl[T]
	options map[string]any
	indexes map[string]*IndexMetainfo
}

var (
	tables = map[reflect.Type]any{}
)

func Table[T any]() *_Table[T] {
	modeltype := reflect.TypeOf((*T)(nil)).Elem()
	mv, ok := tables[modeltype]
	if ok {
		return mv.(*_Table[T])
	}

	tab := &_Table[T]{
		_TypeInfo: gettypeinfo[T](modeltype),
		ddl:       new(_Ddl[T]),
	}
	tables[modeltype] = tab
	tab.ddl.table = tab
	return tab
}

func (table *_Table[T]) ModelPtr() *T {
	return table.modelptr
}

func (tab *_Table[T]) DDL() *_Ddl[T] {
	return tab.ddl
}
