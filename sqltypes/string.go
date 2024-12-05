package sqltypes

type stringType struct {
	typeCommon[string, stringType]
}

func Char(name string, size int) *stringType {
	obj := &stringType{}
	obj.pairs = append(obj.pairs, pair{"name", name})
	obj.pairs = append(obj.pairs, pair{"sqltype", sqlType{"char", []any{size}}})
	return obj
}

func Varchar(name string, size int) *stringType {
	obj := &stringType{}
	obj.pairs = append(obj.pairs, pair{"name", name})
	obj.pairs = append(obj.pairs, pair{"sqltype", sqlType{"varchar", []any{size}}})
	return obj
}

func Text(name string) *stringType {
	obj := &stringType{}
	obj.pairs = append(obj.pairs, pair{"name", name})
	obj.pairs = append(obj.pairs, pair{"sqltype", sqlType{kind: "text"}})
	return obj
}
