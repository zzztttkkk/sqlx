package sqlx

import (
	"fmt"
	"reflect"
	"unsafe"
)

type _Table[T any] struct {
	modeltype reflect.Type
	ptr       *T
	begin     int64
	ddl       *_Ddl[T]
	fields    []_Field
}

var (
	modelptrs = map[reflect.Type]any{}
	tables    = map[reflect.Type]any{}
)

func ModelPtr[T any]() *T {
	var tmp [0]T
	et := reflect.TypeOf(tmp).Elem()

	v, ok := modelptrs[et]
	if ok {
		return v.(*T)
	}
	ptr := new(T)
	modelptrs[et] = ptr
	return ptr
}

func Mptr[T any]() *T {
	return ModelPtr[T]()
}

func Table[T any]() *_Table[T] {
	modeltype := reflect.TypeOf((*T)(nil)).Elem()
	mv, ok := tables[modeltype]
	if ok {
		return mv.(*_Table[T])
	}

	tab := &_Table[T]{
		ptr: ModelPtr[T](),
		ddl: new(_Ddl[T]),
	}
	tables[modeltype] = tab

	tab.modeltype = modeltype
	modelptrs[tab.modeltype] = tab.ptr
	tab.begin = int64(uintptr(unsafe.Pointer(tab.ptr)))
	tab.ddl.table = tab
	tab.init()
	return tab
}

func (table *_Table[T]) ModelPtr() *T {
	return table.ptr
}

func (tab *_Table[T]) init() {
	addfield(&tab.fields, tab.ptr, tab.begin)
}

func addfield(fs *[]_Field, ptr any, begin int64) {
	rv := reflect.ValueOf(ptr).Elem()
	if rv.Kind() != reflect.Struct {
		panic(fmt.Errorf("sqlx: `%s` is not a struct type", rv.Type()))
	}
	rt := rv.Type()

	for i := 0; i < rv.NumField(); i++ {
		fv := rv.Field(i)
		ft := rt.Field(i)
		if !ft.IsExported() {
			continue
		}

		if ft.Anonymous {
			addfield(fs, fv.Addr().Interface(), begin)
			continue
		}

		offset := int64(uintptr(fv.Addr().Pointer())) - begin
		*fs = append(*fs, _Field{
			Offset: offset,
			Field:  ft,
		})
	}
}

func (tab *_Table[T]) fieldbyptr(ptr unsafe.Pointer) *_Field {
	return tab.fieldbyoffset(int64(uintptr(ptr)) - tab.begin)
}

func (tab *_Table[T]) fieldbyoffset(offset int64) *_Field {
	for idx := range tab.fields {
		fp := &tab.fields[idx]
		if fp.Offset == offset {
			return fp
		}
	}
	panic(fmt.Errorf("sqlx: can not find field, %s, %d", tab.modeltype, offset))
}

func (tab *_Table[T]) DDL() *_Ddl[T] {
	return tab.ddl
}