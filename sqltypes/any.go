package sqltypes

import (
	"database/sql"
	"unsafe"
)

type anyTypeBuilder[T any] struct {
	typecommonBuilder[T, anyTypeBuilder[T]]
}

func Type[T any](ptr *T, sqltype string, sqltypeargs ...any) *anyTypeBuilder[T] {
	ins := &anyTypeBuilder[T]{}
	ins.ptr = unsafe.Pointer(ptr)
	ins.sqltype(sqltype, sqltypeargs...)
	return ins
}

func NullableType[T any](ptr *sql.Null[T], sqltype string, sqltypeargs ...any) *anyTypeBuilder[T] {
	return Type((*T)(unsafe.Pointer(ptr)), sqltype, sqltypeargs...).Nullable()
}
