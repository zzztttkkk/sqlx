package sqlx

import (
	"database/sql"
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

type DdlOptions struct {
	Name         string
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
