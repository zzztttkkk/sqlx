package sqltypes

import "unsafe"

type stringTypeBuilder struct {
	typecommonBuilder[string, stringTypeBuilder]
}

func Char(ptr *string, name string, length int) *stringTypeBuilder {
	ins := &stringTypeBuilder{}
	ins.ptr = unsafe.Pointer(ptr)
	ins.sqltype(name, "char", length)
	return ins
}

func Varchar(ptr *string, name string, length int) *stringTypeBuilder {
	ins := &stringTypeBuilder{}
	ins.ptr = unsafe.Pointer(ptr)
	ins.sqltype(name, "varchar", length)
	return ins
}

func Text(ptr *string, name string) *stringTypeBuilder {
	ins := &stringTypeBuilder{}
	ins.ptr = unsafe.Pointer(ptr)
	ins.sqltype(name, "text")
	return ins
}
