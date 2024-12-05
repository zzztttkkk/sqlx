package sqlx

import "reflect"

type SqlTable struct {
	Name    string
	GoType  reflect.Type
	Fields  []*SqlField
	Indexes []*Index
	Options map[string]string
}
