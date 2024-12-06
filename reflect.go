package sqlx

import (
	"fmt"
	"reflect"
)

type TypeInfo struct {
	GoType reflect.Type
	Fields []*FieldMetaInfo
}

var (
	typeinfos  = map[reflect.Type]*TypeInfo{}
	fieldinfos = map[reflect.Type]*FieldMetaInfo{}
	tables     = map[reflect.Type]*TableMetaInfo{}
)

func QuerySqlField(fieldtype reflect.Type) *FieldMetaInfo {
	return fieldinfos[fieldtype]
}

func mkTypeInfo(type_ reflect.Type, checkTableType bool) *TypeInfo {
	if type_.Kind() != reflect.Struct {
		panic(fmt.Errorf("sqlx: expected a struct type, but got `%s`", type_))
	}

	v, ok := typeinfos[type_]
	if ok {
		return v
	}

	v = &TypeInfo{GoType: type_}
	for i := 0; i < type_.NumField(); i++ {
		ft := type_.Field(i)
		if ft.Anonymous {
			tinfo := mkTypeInfo(ft.Type, false)
			for _, fv := range tinfo.Fields {
				var ptr = &FieldMetaInfo{}
				*ptr = *fv
				v.Fields = append(v.Fields, ptr)
			}
			continue
		}
		if !ft.IsExported() || !ft.Type.Implements(typeofIfaceField) {
			continue
		}

		fv := reflect.New(ft.Type).Elem().Interface().(ifaceField)
		metatype := fv.__sqlxfield__metatype()
		tabletype := fv.__sqlxfield__tabletype()

		if checkTableType && tabletype != type_ {
			panic(fmt.Errorf("sqlx: `%s.%s`'s table type is wrong", type_, ft.Name))
		}
		sf := reflect.New(metatype).Elem().Interface().(IFieldMeta).FieldMetaInfo()
		sf.metaType = metatype
		sf.fieldType = ft.Type
		v.Fields = append(v.Fields, &sf)
	}
	for _, field := range v.Fields {
		field.structType = type_
		fieldinfos[field.fieldType] = field
	}
	typeinfos[type_] = v
	return v
}

func mkTableMetaInfo(typeinfo *TypeInfo) *TableMetaInfo {
	tmi := reflect.New(typeinfo.GoType).Elem().Interface().(ITable).TableMetaInfo()
	tmi.goType = typeinfo.GoType
	tables[tmi.goType] = &tmi
	return &tmi
}

func RegisterTable(val ITable) *TableMetaInfo {
	return mkTableMetaInfo(mkTypeInfo(reflect.TypeOf(val), true))
}

func RegisterScanableStructs(vals ...any) {
	for _, val := range vals {
		vt := reflect.TypeOf(val)
		if vt.Kind() != reflect.Struct {
			continue
		}
		mkTypeInfo(vt, false)
	}
}
