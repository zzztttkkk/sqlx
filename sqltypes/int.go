package sqltypes

import (
	"database/sql"
	"reflect"
	"strings"
	"unsafe"
)

type ints interface {
	~int8 | ~int16 | ~int32 | ~int64 | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~int | ~uint
}

type intTypeBuilder[T ints] struct {
	typecommonBuilder[T, intTypeBuilder[T]]
}

func (builder *intTypeBuilder[T]) AutoIncr() *intTypeBuilder[T] {
	builder.pairs = append(builder.pairs, pair{"autoincr", true})
	return builder
}

func getintkind[T ints]() (string, int) {
	tt := reflect.TypeOf(T(0))
	if strings.HasPrefix(tt.Name(), "u") {
		return "uint", tt.Bits()
	}
	return "int", tt.Bits()
}

func Int[T ints](ptr *T) *intTypeBuilder[T] {
	ins := &intTypeBuilder[T]{}
	ins.ptr = unsafe.Pointer(ptr)
	ik, is := getintkind[T]()
	ins.sqltype(ik, is)
	return ins
}

func NullableInt[T ints](ptr *sql.Null[T]) *intTypeBuilder[T] {
	return Int((*T)(unsafe.Pointer(ptr))).Nullable()
}
