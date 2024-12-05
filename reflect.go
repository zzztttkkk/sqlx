package sqlx

import (
	"reflect"
)

type TypeInfo struct {
	Fields []*SqlField
}

type fieldinfosKeyType struct {
	tablletype reflect.Type
	metatype   reflect.Type
}

var (
	typeinfos  = map[reflect.Type]*TypeInfo{}
	fieldinfos = map[fieldinfosKeyType]*SqlField{}
)

func QuerySqlField(tabletype reflect.Type, metatype reflect.Type) *SqlField {
	return fieldinfos[fieldinfosKeyType{tabletype, metatype}]
}

func mktypeinfo(type_ reflect.Type) *TypeInfo {
	v, ok := typeinfos[type_]
	if ok {
		return v
	}

	v = &TypeInfo{}

	for i := 0; i < type_.NumField(); i++ {
		ft := type_.Field(i)
		if !ft.IsExported() {
			continue
		}
		if ft.Anonymous {
			tinfo := mktypeinfo(ft.Type)
			for _, fv := range tinfo.Fields {
				var ptr = &SqlField{}
				*ptr = *fv
				v.Fields = append(v.Fields, ptr)
			}
			continue
		}
		if !ft.Type.Implements(typeofIfaceField) {
			continue
		}

		metatype := reflect.New(ft.Type).Elem().Interface().(ifaceField).__sqlxfield__metatype()
		sf := reflect.New(metatype).Elem().Interface().(IFieldMeta).SqlField()
		sf.metaType = metatype
		v.Fields = append(v.Fields, &sf)
	}
	for _, field := range v.Fields {
		field.structType = type_
		fieldinfos[fieldinfosKeyType{type_, field.metaType}] = field
	}
	typeinfos[type_] = v
	return v
}
