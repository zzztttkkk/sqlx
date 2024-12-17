package sqlx

import (
	"reflect"

	"github.com/zzztttkkk/reflectx"
)

type _Type[T any] struct {
	*reflectx.TypeInfo[DdlOptions]
	scheme  *_Scheme[T]
	options map[string]any
	indexes map[string]*IndexMetainfo
}

var (
	types = map[reflect.Type]any{}
)

func Type[T any]() *_Type[T] {
	modeltype := reflectx.Typeof[T]()
	mv, ok := types[modeltype]
	if ok {
		return mv.(*_Type[T])
	}

	tab := &_Type[T]{
		TypeInfo: reflectx.TypeInfoOf[T, DdlOptions](),
		scheme:   new(_Scheme[T]),
	}
	types[modeltype] = tab
	tab.scheme.table = tab
	return tab
}

func (tab *_Type[T]) Scheme() *_Scheme[T] {
	return tab.scheme
}
