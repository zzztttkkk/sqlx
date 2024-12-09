package sqlx

import (
	"database/sql"
	"reflect"
)

type IndexMetainfo struct {
	Name string
}

type FieldMetainfo struct {
	Name       string
	SqlType    string
	PrimaryKey bool
	Unique     bool
	Nullable   bool
	Check      string
	Default    sql.NullString
	Comment    string
}

type _Field struct {
	Offset   int64
	Field    reflect.StructField
	Metainfo *FieldMetainfo
}
