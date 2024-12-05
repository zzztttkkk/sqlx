package sqltypes

import "unsafe"

type pair struct {
	key string
	val any
}

type sqlType struct {
	kind string
	args []any
}

type typeCommon[T any, Self any] struct {
	pairs []pair
	opts  map[string]string
}

func (tc *typeCommon[T, Self]) self() *Self {
	return (*Self)(unsafe.Pointer(tc))
}

func (tc *typeCommon[T, Self]) Primary() *Self {
	tc.pairs = append(tc.pairs, pair{"primary", true})
	return tc.self()
}

func (tc *typeCommon[T, Self]) Unique() *Self {
	tc.pairs = append(tc.pairs, pair{"unique", true})
	return tc.self()
}

func (tc *typeCommon[T, Self]) Default(dv T) *Self {
	tc.pairs = append(tc.pairs, pair{"default", dv})
	return tc.self()
}

func (tc *typeCommon[T, Self]) DefaultExpr(expr string) *Self {
	tc.pairs = append(tc.pairs, pair{"defaultexpr", expr})
	return tc.self()
}

func (tc *typeCommon[T, Self]) Nullable() *Self {
	tc.pairs = append(tc.pairs, pair{"nullable", true})
	return tc.self()
}

func (tc *typeCommon[T, Self]) Comment(comment string) *Self {
	tc.pairs = append(tc.pairs, pair{"comment", comment})
	return tc.self()
}

func (tc *typeCommon[T, Self]) Check(check string) *Self {
	tc.pairs = append(tc.pairs, pair{"check", check})
	return tc.self()
}

func (tc *typeCommon[T, Self]) Options(k string, v string) *Self {
	if tc.opts == nil {
		tc.opts = map[string]string{}
	}
	tc.opts[k] = v
	return tc.self()
}
