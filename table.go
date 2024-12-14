package sqlx

import (
	"reflect"

	"github.com/zzztttkkk/reflectx"
)

type _Table[T any] struct {
	*reflectx.TypeInfo[DdlOptions]
	scheme  *_Scheme[T]
	options map[string]any
	indexes map[string]*IndexMetainfo
}

var (
	tables = map[reflect.Type]any{}
)

func Table[T any]() *_Table[T] {
	modeltype := reflectx.Typeof[T]()
	mv, ok := tables[modeltype]
	if ok {
		return mv.(*_Table[T])
	}

	tab := &_Table[T]{
		TypeInfo: reflectx.TypeInfoOf[T, DdlOptions](),
		scheme:   new(_Scheme[T]),
	}
	tables[modeltype] = tab
	tab.scheme.table = tab
	return tab
}

func (tab *_Table[T]) Scheme() *_Scheme[T] {
	return tab.scheme
}
