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

type FieldMetainfo struct {
	Name        string
	SqlType     string
	SqlTypeArgs []any
	PrimaryKey  bool
	Unique      bool
	Nullable    bool
	Check       string
	Default     sql.NullString
	Comment     string
	AutoIncr    bool
}

type _Field struct {
	Offset   int64
	Field    reflect.StructField
	Metainfo *FieldMetainfo
}
