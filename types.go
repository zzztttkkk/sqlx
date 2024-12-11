package sqlx

import (
	"database/sql"
	"reflect"
)

type IndexField struct {
	Name    string
	Order   OrderKind
	Options map[string]any
}

type OrderKind int

const (
	OrderKindAuto = OrderKind(iota)
	OrderKindAsc
	OrderKindDesc
)

type IndexMetainfo struct {
	Name    string
	Unique  bool
	Fields  []IndexField
	Options map[string]any
}

type FieldDdlOptions struct {
	SqlType      string
	SqlTypeArgs  []any
	PrimaryKey   bool
	Unique       bool
	Nullable     bool
	CheckExpr    string
	DefaultExpr  sql.Null[string]
	DefaultValue sql.Null[any]
	Comment      string
	AutoIncr     bool
}

type FieldMetainfo struct {
	Name       string
	DdlOptions FieldDdlOptions
}

type _Field struct {
	Offset    int64
	Field     reflect.StructField
	Metainfo  *FieldMetainfo
}

func (f *_Field) setmeta(meta *FieldMetainfo) {
	if f.Metainfo == nil {
		f.Metainfo = meta
		return
	}
	if meta.Name != "" {
		f.Metainfo.Name = meta.Name
	}
	f.Metainfo.DdlOptions = meta.DdlOptions
}
