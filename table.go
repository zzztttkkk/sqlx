package sqlx

import "reflect"

type TableMetaInfo struct {
	Name    string
	goType  reflect.Type
	fields  []*FieldMetaInfo
	indexes []*Index
	Options map[string]string
}

type ITable interface {
	TableMetaInfo() TableMetaInfo
}
