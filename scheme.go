package sqlx

import (
	"unsafe"

	"github.com/zzztttkkk/lion"
)

type _Scheme[T any] struct {
	table *_Type[T]
	name  string
}

type ischeme interface {
	updatefield(field SqlField)
}

func (scheme *_Scheme[T]) updatefield(field SqlField) {
	fp := scheme.table.FieldByUnsafePtr(field.Ptr)
	fp.UpdateMetainfo(field.Metainfo)
}

func (scheme *_Scheme[T]) Field(field SqlField) *_Scheme[T] {
	scheme.updatefield(field)
	return scheme
}

func (scheme *_Scheme[T]) Tablename(name string) *_Scheme[T] {
	scheme.name = name
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
	Fields  []*DdlOptions
	Indexes []*IndexMetainfo
	Options map[string]any
}

func (scheme *_Scheme[T]) Finish() *TableMetainfo {
	var nfs []lion.Field[DdlOptions]
	for _, f := range scheme.table.Fields {
		if f.Metainfo() == nil {
			continue
		}
		nfs = append(nfs, f)
	}
	scheme.table.Fields = nfs
	tmi := &TableMetainfo{
		Name:    scheme.name,
		Options: scheme.table.options,
	}
	for _, f := range scheme.table.Fields {
		mi := f.Metainfo()
		mi.Name = f.Name()
		tmi.Fields = append(tmi.Fields, mi)
	}
	for _, idx := range scheme.table.indexes {
		tmi.Indexes = append(tmi.Indexes, idx)
	}
	return tmi
}

type SqlField struct {
	Ptr      unsafe.Pointer
	Metainfo *DdlOptions
}

func (sf SqlField) AddToScheme(scheme ischeme) {
	scheme.updatefield(sf)
}
