package sqlx

import (
	"database/sql"
	"fmt"
	"reflect"
)

type DdlOptions struct {
	SqlType  string
	Primary  bool
	Unique   bool
	Default  sql.NullString
	Nullable bool
	Comment  string
	AutoIncr bool
	Check    string
	Options  map[string]string
}

type FieldMetaInfo struct {
	Name   string
	GoType reflect.Type
	Table  *TableMetaInfo
	Ddl    DdlOptions

	structType reflect.Type
	metaType   reflect.Type
	fieldType  reflect.Type
}

type IFieldMeta interface {
	FieldMetaInfo() FieldMetaInfo
}

type Field[V any, Meta IFieldMeta, TablePtr ITable] struct {
	Value V
	meta  [0]Meta
	table [0]TablePtr
}

func (field Field[V, Meta, TablePtr]) __sqlxfield__tabletype() reflect.Type {
	tet := reflect.TypeOf(field.table).Elem()
	if tet.Kind() != reflect.Pointer {
		panic(fmt.Errorf("sqlx: bad table type, should be the pointer of struct"))
	}
	return tet.Elem()
}

func (field Field[V, M, T]) __sqlxfield__metatype() reflect.Type {
	met := reflect.TypeOf(field.meta).Elem()
	if met.Kind() == reflect.Pointer {
		met = met.Elem()
	}
	return met
}

func (field *Field[V, M, T]) SqlField() *FieldMetaInfo {
	return fieldinfos[reflect.TypeOf(field).Elem()]
}

func (field *Field[V, M, T]) ScanPtr() any {
	return &field.Value
}

type ISqlField interface {
	SqlField() *FieldMetaInfo
	ScanPtr() any
}

type ifaceField interface {
	__sqlxfield__metatype() reflect.Type
	__sqlxfield__tabletype() reflect.Type
}

var typeofIfaceField = reflect.TypeOf((*ifaceField)(nil)).Elem()

type nonFieldMeta struct{}

func (n nonFieldMeta) FieldMetaInfo() FieldMetaInfo {
	return FieldMetaInfo{}
}

type nonTableMeta struct{}

func (_ nonTableMeta) TableMetaInfo() TableMetaInfo {
	return TableMetaInfo{}
}

var _ ifaceField = Field[int, nonFieldMeta, nonTableMeta]{}
