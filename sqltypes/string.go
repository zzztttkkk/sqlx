package sqltypes

import (
	"database/sql"
	"unsafe"
)

type stringTypeBuilder struct {
	typecommonBuilder[string, stringTypeBuilder]
}

func Char(ptr *string, name string, length int) *stringTypeBuilder {
	ins := &stringTypeBuilder{}
	ins.ptr = unsafe.Pointer(ptr)
	ins.sqltype(name, "char", length)
	return ins
}

func NullableChar(ptr *sql.Null[string], name string, length int) *stringTypeBuilder {
	return Char((*string)(unsafe.Pointer(ptr)), name, length).Nullable()
}

func Varchar(ptr *string, name string, length int) *stringTypeBuilder {
	ins := &stringTypeBuilder{}
	ins.ptr = unsafe.Pointer(ptr)
	ins.sqltype(name, "varchar", length)
	return ins
}

func NullableVarchar(ptr *sql.Null[string], name string, length int) *stringTypeBuilder {
	return Varchar((*string)(unsafe.Pointer(ptr)), name, length).Nullable()
}

func Text(ptr *string, name string) *stringTypeBuilder {
	ins := &stringTypeBuilder{}
	ins.ptr = unsafe.Pointer(ptr)
	ins.sqltype(name, "text")
	return ins
}

func NullableText(ptr *sql.Null[string], name string) *stringTypeBuilder {
	return Text((*string)(unsafe.Pointer(ptr)), name).Nullable()
}
