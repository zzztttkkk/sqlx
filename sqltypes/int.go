package sqltypes

import "unsafe"

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

func TinyInt(ptr *int8, name string) *intTypeBuilder[int8] {
	ins := &intTypeBuilder[int8]{}
	ins.ptr = unsafe.Pointer(ptr)
	ins.sqltype(name, "int", 8)
	return ins
}

func SmallInt(ptr *int16, name string) *intTypeBuilder[int16] {
	ins := &intTypeBuilder[int16]{}
	ins.ptr = unsafe.Pointer(ptr)
	ins.sqltype(name, "int", 16)
	return ins
}

func Int(ptr *int32, name string) *intTypeBuilder[int32] {
	ins := &intTypeBuilder[int32]{}
	ins.ptr = unsafe.Pointer(ptr)
	ins.sqltype(name, "int", 32)
	return ins
}

func BigInt(ptr *int64, name string) *intTypeBuilder[int64] {
	ins := &intTypeBuilder[int64]{}
	ins.ptr = unsafe.Pointer(ptr)
	ins.sqltype(name, "int", 64)
	return ins
}

func TinyUint(ptr *uint8, name string) *intTypeBuilder[uint8] {
	ins := &intTypeBuilder[uint8]{}
	ins.ptr = unsafe.Pointer(ptr)
	ins.sqltype(name, "uint", 8)
	return ins
}

func SmallUint(ptr *uint16, name string) *intTypeBuilder[uint16] {
	ins := &intTypeBuilder[uint16]{}
	ins.ptr = unsafe.Pointer(ptr)
	ins.sqltype(name, "uint", 16)
	return ins
}

func Uint(ptr *uint32, name string) *intTypeBuilder[uint32] {
	ins := &intTypeBuilder[uint32]{}
	ins.ptr = unsafe.Pointer(ptr)
	ins.sqltype(name, "uint", 32)
	return ins
}

func BigUint(ptr *uint64, name string) *intTypeBuilder[uint64] {
	ins := &intTypeBuilder[uint64]{}
	ins.ptr = unsafe.Pointer(ptr)
	ins.sqltype(name, "uint", 64)
	return ins
}

func CastPtr[T any](v unsafe.Pointer) *T {
	return (*T)(v)
}
