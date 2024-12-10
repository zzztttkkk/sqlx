package sqlx

import (
	"fmt"
	"reflect"
	"strings"
	"unsafe"
)

type _TypeInfo[T any] struct {
	modeltype reflect.Type
	modelptr  *T
	ptrnum    int64
	fields    []_Field
}

func (tab *_TypeInfo[T]) fieldbyptr(ptr unsafe.Pointer) (*_Field, bool) {
	return tab.fieldbyoffset(int64(uintptr(ptr)) - tab.ptrnum)
}

func (tab *_TypeInfo[T]) mustfieldbyptr(ptr unsafe.Pointer) *_Field {
	return tab.mustfieldbyoffset(int64(uintptr(ptr)) - tab.ptrnum)
}

func (tab *_TypeInfo[T]) fieldbyoffset(offset int64) (*_Field, bool) {
	for idx := range tab.fields {
		fp := &tab.fields[idx]
		if fp.Offset == offset {
			return fp, true
		}
	}
	return nil, false
}

func (tab *_TypeInfo[T]) mustfieldbyoffset(offset int64) *_Field {
	ptr, ok := tab.fieldbyoffset(offset)
	if ok {
		return ptr
	}
	panic(fmt.Errorf("sqlx: can not find field, %s, %d", tab.modeltype, offset))
}

func gettag(field *reflect.StructField, tags ...string) string {
	for _, tag := range tags {
		v := field.Tag.Get(tag)
		if v != "" {
			return v
		}
	}
	return ""
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
		dbtag := gettag(&ft, "sqlx", "db", "sql")
		if dbtag == "-" || !ft.IsExported() {
			continue
		}
		if ft.Anonymous {
			addfield(fs, fv.Addr().Interface(), begin)
			continue
		}

		offset := int64(uintptr(fv.Addr().Pointer())) - begin

		f := _Field{
			Offset:   offset,
			Field:    ft,
			Metainfo: &FieldMetainfo{},
		}

		parts := strings.Split(dbtag, ",")
		name := strings.TrimSpace(parts[0])
		if name == "" {
			name = ft.Name
		}
		f.Metainfo.Name = name
		*fs = append(*fs, f)
	}
}

var (
	modelptrs = map[reflect.Type]any{}
	typeinfos = map[reflect.Type]any{}
)

func Mptr[T any]() *T {
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

func ModelPtr[T any]() *T {
	return Mptr[T]()
}

func gettypeinfo[T any](tt reflect.Type) *_TypeInfo[T] {
	if tt == nil {
		tt = reflect.TypeOf((*T)(nil)).Elem()
	}
	v, ok := typeinfos[tt]
	if ok {
		return v.(*_TypeInfo[T])
	}

	ti := &_TypeInfo[T]{
		modeltype: tt,
		modelptr:  Mptr[T](),
	}
	ti.ptrnum = int64(uintptr(unsafe.Pointer(ti.modelptr)))
	typeinfos[tt] = ti
	addfield(&ti.fields, ti.modelptr, ti.ptrnum)
	return ti
}

func setfieldmeta[T any](fptr unsafe.Pointer, meta *FieldMetainfo) {
	ti := gettypeinfo[T](nil)
	fv := ti.mustfieldbyptr(fptr)
	fv.setmeta(meta)
}
