package sqltypes

type IntType interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

type intType[T IntType] struct {
	typeCommon[T, intType[T]]
}

func TinyInt(name string) *intType[int8] {
	obj := &intType[int8]{}
	obj.pairs = append(obj.pairs, pair{"name", name})
	obj.pairs = append(obj.pairs, pair{"sqltype", sqlType{kind: "int", args: []any{8}}})
	return obj
}

func SmallInt(name string) *intType[int16] {
	obj := &intType[int16]{}
	obj.pairs = append(obj.pairs, pair{"name", name})
	obj.pairs = append(obj.pairs, pair{"sqltype", sqlType{kind: "int", args: []any{16}}})
	return obj
}

func Int(name string) *intType[int32] {
	obj := &intType[int32]{}
	obj.pairs = append(obj.pairs, pair{"name", name})
	obj.pairs = append(obj.pairs, pair{"sqltype", sqlType{kind: "int", args: []any{32}}})
	return obj
}

func BigInt(name string) *intType[int64] {
	obj := &intType[int64]{}
	obj.pairs = append(obj.pairs, pair{"name", name})
	obj.pairs = append(obj.pairs, pair{"sqltype", sqlType{kind: "int", args: []any{64}}})
	return obj
}

func TinyUint(name string) *intType[uint8] {
	obj := &intType[uint8]{}
	obj.pairs = append(obj.pairs, pair{"name", name})
	obj.pairs = append(obj.pairs, pair{"sqltype", sqlType{kind: "uint", args: []any{8}}})
	return obj
}

func SmallUint(name string) *intType[uint16] {
	obj := &intType[uint16]{}
	obj.pairs = append(obj.pairs, pair{"name", name})
	obj.pairs = append(obj.pairs, pair{"sqltype", sqlType{kind: "uint", args: []any{16}}})
	return obj
}

func Uint(name string) *intType[uint32] {
	obj := &intType[uint32]{}
	obj.pairs = append(obj.pairs, pair{"name", name})
	obj.pairs = append(obj.pairs, pair{"sqltype", sqlType{kind: "uint", args: []any{32}}})
	return obj
}

func BigUint(name string) *intType[uint64] {
	obj := &intType[uint64]{}
	obj.pairs = append(obj.pairs, pair{"name", name})
	obj.pairs = append(obj.pairs, pair{"sqltype", sqlType{kind: "uint", args: []any{64}}})
	return obj
}
