package sqltypes

import (
	"database/sql"
	"unsafe"
)

type boolTypeBuilder struct {
	typecommonBuilder[bool, boolTypeBuilder]
}

func Bool(ptr *bool) *boolTypeBuilder {
	ins := &boolTypeBuilder{}
	ins.ptr = unsafe.Pointer(ptr)
	ins.sqltype("bool")
	return ins
}

func NullableBool(ptr *sql.Null[bool]) *boolTypeBuilder {
	return Bool((*bool)(unsafe.Pointer(ptr))).Nullable()
}
