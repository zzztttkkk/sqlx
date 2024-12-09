package sqltypes

import (
	"database/sql"
	"unsafe"
)

type boolTypeBuilder struct {
	typecommonBuilder[bool, boolTypeBuilder]
}

func Bool(ptr *bool, name string) *boolTypeBuilder {
	ins := &boolTypeBuilder{}
	ins.ptr = unsafe.Pointer(ptr)
	ins.sqltype(name, "bool")
	return ins
}

func NullableBool(ptr *sql.Null[bool], name string) *boolTypeBuilder {
	return Bool((*bool)(unsafe.Pointer(ptr)), name).Nullable()
}
