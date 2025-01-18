package sqlx

import (
	"reflect"

	"github.com/zzztttkkk/lion"
)

type _Type[T any] struct {
	*lion.TypeInfo
	scheme  *_Scheme[T]
	options map[string]any
	indexes map[string]*IndexMetainfo
}

var (
	types = map[reflect.Type]any{}
)

func Type[T any]() *_Type[T] {
	modeltype := lion.Typeof[T]()
	mv, ok := types[modeltype]
	if ok {
		return mv.(*_Type[T])
	}

	tab := &_Type[T]{
		TypeInfo: lion.TypeInfoOf[T](),
		scheme:   new(_Scheme[T]),
	}
	types[modeltype] = tab
	tab.scheme.table = tab
	return tab
}

func (tab *_Type[T]) Scheme() *_Scheme[T] {
	return tab.scheme
}
