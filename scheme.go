package sqlx

import (
	"fmt"
	"reflect"
	"unsafe"
)

type _Scheme[T any] struct {
	table *_Table[T]
	name  string
}

type ischeme interface {
	addfield(field SqlField)
}

func (scheme *_Scheme[T]) addfield(field SqlField) {
	fp := scheme.table.mustfieldbyptr(field.Ptr)
	fp.setmeta(field.Metainfo)
}

func (scheme *_Scheme[T]) Field(field SqlField) *_Scheme[T] {
	scheme.addfield(field)
	return scheme
}

func (scheme *_Scheme[T]) Tablename(name string) *_Scheme[T] {
	scheme.name = name
	return scheme
}

func (scheme *_Scheme[T]) Mixed(ptr IMixed) *_Scheme[T] {
	pt := reflect.TypeOf(ptr).Elem()
	if pt.Kind() != reflect.Struct {
		panic(fmt.Errorf("sqlx: mixed is not a struct type, %s", pt))
	}
	mt := scheme.table.modeltype
	mp := scheme.table.modelptr
	mbegin := scheme.table.ptrnum
	mv := reflect.ValueOf(mp).Elem()

	idx := -1
	for i := 0; i < mt.NumField(); i++ {
		ft := mt.Field(i)
		if ft.Anonymous && ft.Type == pt {
			if idx < 0 {
				idx = i
				continue
			}
			panic(fmt.Errorf("sqlx: repeatedly mix same type, %s, %s", scheme.table.modeltype, pt))
		}
		continue
	}
	if idx < 0 {
		panic(fmt.Errorf("sqlx: has no mix type, %s, %s", scheme.table.modeltype, pt))
	}
	mmoffset := int64(uintptr(mv.Field(idx).Addr().Pointer())) - mbegin

	var begin = int64(uintptr(reflect.ValueOf(ptr).Pointer()))
	for _, mf := range ptr.MixedFields() {
		offset := int64(uintptr(mf.Ptr)) - begin + mmoffset
		fp := scheme.table.mustfieldbyoffset(offset)
		fp.setmeta(mf.Metainfo)
	}
	return scheme
}

func (scheme *_Scheme[T]) Option(k string, v any) *_Scheme[T] {
	if scheme.table.options == nil {
		scheme.table.options = map[string]any{}
	}
	scheme.table.options[k] = v
	return scheme
}

type TableMetainfo struct {
	Name    string
	Fields  []*FieldMetainfo
	Indexes []*IndexMetainfo
	Options map[string]any
}

func (scheme *_Scheme[T]) Finish() *TableMetainfo {
	var nfs []_Field
	for _, f := range scheme.table.fields {
		if f.Metainfo == nil {
			continue
		}
		nfs = append(nfs, f)
	}
	scheme.table.fields = nfs
	tmi := &TableMetainfo{
		Name:    scheme.name,
		Options: scheme.table.options,
	}
	for _, f := range scheme.table.fields {
		tmi.Fields = append(tmi.Fields, f.Metainfo)
	}
	for _, idx := range scheme.table.indexes {
		tmi.Indexes = append(tmi.Indexes, idx)
	}
	return tmi
}

type SqlField struct {
	Ptr      unsafe.Pointer
	Metainfo *FieldMetainfo
}

func (sf SqlField) AddToScheme(scheme ischeme) {
	scheme.addfield(sf)
}

type IMixed interface {
	MixedFields() []SqlField
}