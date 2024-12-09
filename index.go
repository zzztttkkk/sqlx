package sqlx

import (
	"fmt"
	"unsafe"
)

type indexBuilder[T any] struct {
	table *_Table[T]
	meta  *IndexMetainfo
}

func Index[T any](name string) *indexBuilder[T] {
	ins := &indexBuilder[T]{
		table: Table[T](),
		meta:  &IndexMetainfo{Name: name},
	}
	if ins.table.indexes == nil {
		ins.table.indexes = map[string]*IndexMetainfo{}
	}
	ins.table.indexes[name] = ins.meta
	return ins
}

func (builder *indexBuilder[T]) Unique() *indexBuilder[T] {
	builder.meta.Unique = true
	return builder
}

func (builder *indexBuilder[T]) Option(k string, val any) *indexBuilder[T] {
	if builder.meta.Options == nil {
		builder.meta.Options = map[string]any{}
	}
	builder.meta.Options[k] = val
	return builder
}

func (build *indexBuilder[T]) Field(ptr unsafe.Pointer, order OrderKind, opts map[string]any) *indexBuilder[T] {
	field, ok := build.table.fieldbyptr(ptr)
	if !ok {
		panic(fmt.Errorf(
			"sqlx: failed to get field metainfo through pointer when creating index. Did you set all the fields ? or pass a wrong pointer, TableType(%s), IndexName(%s)",
			build.table.modeltype,
			build.meta.Name,
		))
	}
	if field.Metainfo == nil {
		panic(fmt.Errorf("sqlx: field metainfo is nil, TableType(%s), FieldName(%s)", build.table.modeltype, field.Field.Name))
	}
	build.meta.Fields = append(build.meta.Fields, IndexField{
		Name:    field.Metainfo.Name,
		Order:   order,
		Options: opts,
	})
	return build
}
