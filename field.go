package sqlx

import (
	"database/sql"
	"reflect"
)

type SqlField struct {
	Name       string
	SqlType    string
	GoType     reflect.Type
	structType reflect.Type
	metaType   reflect.Type
	fieldType  reflect.Type

	Primary  bool
	Unique   bool
	Default  sql.NullString
	Nullable bool
	Comment  string
	AutoIncr bool
	Check    string
	Options  map[string]string

	Table *SqlTable
}

type IFieldMeta interface {
	SqlField() SqlField
}

type Field[V any, Meta IFieldMeta, TablePtr any] struct {
	Value V

	meta [0]Meta
	_    [0]TablePtr
}

func (field Field[V, M, T]) __sqlxfield__metatype() reflect.Type {
	return reflect.TypeOf(field.meta).Elem()
}

func (field *Field[V, M, T]) SqlField() *SqlField {
	return fieldinfos[reflect.TypeOf(field).Elem()]
}

type ifaceField interface {
	__sqlxfield__metatype() reflect.Type
}

var typeofIfaceField = reflect.TypeOf((*ifaceField)(nil)).Elem()

type nonFieldMeta struct{}

func (n nonFieldMeta) SqlField() SqlField {
	panic("unimplemented")
}

var _ IFieldMeta = nonFieldMeta{}

var _ ifaceField = Field[int, nonFieldMeta, any]{}
