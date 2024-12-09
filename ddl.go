package sqlx

import (
	"fmt"
	"reflect"
	"unsafe"
)

type _Ddl[T any] struct {
	table *_Table[T]
	name  string
}

type iddl interface {
	addfield(field SqlField)
}

func (ddl *_Ddl[T]) addfield(field SqlField) {
	fp := ddl.table.fieldbyptr(field.Ptr)
	fp.Metainfo = field.Metainfo
}

func (ddl *_Ddl[T]) Field(field SqlField) *_Ddl[T] {
	ddl.addfield(field)
	return ddl
}

func (ddl *_Ddl[T]) Tablename(name string) *_Ddl[T] {
	ddl.name = name
	return ddl
}

func (ddl *_Ddl[T]) Mixed(ptr IMixed) *_Ddl[T] {
	pt := reflect.TypeOf(ptr).Elem()
	if pt.Kind() != reflect.Struct {
		panic(fmt.Errorf("sqlx: mixed is not a struct type, %s", pt))
	}
	mt := ddl.table.modeltype
	mp := ddl.table.ptr
	mbegin := ddl.table.begin
	mv := reflect.ValueOf(mp).Elem()

	idx := -1
	for i := 0; i < mt.NumField(); i++ {
		ft := mt.Field(i)
		if ft.Anonymous && ft.Type == pt {
			if idx < 0 {
				idx = i
				continue
			}
			panic(fmt.Errorf("sqlx: repeatedly mix same type, %s, %s", ddl.table.modeltype, pt))
		}
		continue
	}
	if idx < 0 {
		panic(fmt.Errorf("sqlx: has no mix type, %s, %s", ddl.table.modeltype, pt))
	}
	mmoffset := int64(uintptr(mv.Field(idx).Addr().Pointer())) - mbegin

	var begin = int64(uintptr(reflect.ValueOf(ptr).Pointer()))
	for _, mf := range ptr.MixedFields() {
		offset := int64(uintptr(mf.Ptr)) - begin + mmoffset
		fp := ddl.table.fieldbyoffset(offset)
		fp.Metainfo = mf.Metainfo
	}
	return ddl
}

func (ddl *_Ddl[T]) Finish() *_Table[T] {
	var nfs []_Field
	for _, v := range ddl.table.fields {
		if v.Metainfo == nil {
			continue
		}
		nfs = append(nfs, v)
	}
	ddl.table.fields = nfs
	return ddl.table
}

type SqlField struct {
	Ptr      unsafe.Pointer
	Metainfo *FieldMetainfo
}

func (sf SqlField) AddToDdl(ddl iddl) {
	ddl.addfield(sf)
}

type IMixed interface {
	MixedFields() []SqlField
}
