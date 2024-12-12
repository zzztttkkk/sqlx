package sqlx

import (
	"reflect"
)

type _Table[T any] struct {
	*_TypeInfo[T]

	scheme  *_Scheme[T]
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
		scheme:    new(_Scheme[T]),
	}
	tables[modeltype] = tab
	tab.scheme.table = tab
	return tab
}

func (table *_Table[T]) ModelPtr() *T {
	return table.modelptr
}

func (tab *_Table[T]) Scheme() *_Scheme[T] {
	return tab.scheme
}
