package sqltypes

import (
	"database/sql"
	"unsafe"
)

type stringTypeBuilder struct {
	typecommonBuilder[string, stringTypeBuilder]
}

func Char(ptr *string, length int) *stringTypeBuilder {
	ins := &stringTypeBuilder{}
	ins.ptr = unsafe.Pointer(ptr)
	ins.sqltype("char", length)
	return ins
}

func NullableChar(ptr *sql.Null[string], length int) *stringTypeBuilder {
	return Char((*string)(unsafe.Pointer(ptr)), length).Nullable()
}

func Varchar(ptr *string, length int) *stringTypeBuilder {
	ins := &stringTypeBuilder{}
	ins.ptr = unsafe.Pointer(ptr)
	ins.sqltype("varchar", length)
	return ins
}

func NullableVarchar(ptr *sql.Null[string], length int) *stringTypeBuilder {
	return Varchar((*string)(unsafe.Pointer(ptr)), length).Nullable()
}

func Text(ptr *string) *stringTypeBuilder {
	ins := &stringTypeBuilder{}
	ins.ptr = unsafe.Pointer(ptr)
	ins.sqltype("text")
	return ins
}

func NullableText(ptr *sql.Null[string]) *stringTypeBuilder {
	return Text((*string)(unsafe.Pointer(ptr))).Nullable()
}
