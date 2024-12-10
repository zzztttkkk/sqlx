package sqltypes

import (
	"fmt"
	"reflect"
	"unsafe"

	"github.com/zzztttkkk/sqlx"
)

type pair struct {
	key string
	val any
}

type typecommonBuilder[T any, S any] struct {
	ptr   unsafe.Pointer
	pairs []pair
}

type sqltype struct {
	kind string
	args []any
}

func (builder *typecommonBuilder[T, S]) self() *S {
	return (*S)(unsafe.Pointer(builder))
}

func (builder *typecommonBuilder[T, S]) sqltype(sqlkind string, args ...any) *S {
	builder.pairs = append(builder.pairs, pair{"sqltype", sqltype{kind: sqlkind, args: args}})
	return builder.self()
}

func (builder *typecommonBuilder[T, S]) Unique() *S {
	builder.pairs = append(builder.pairs, pair{"unique", true})
	return builder.self()
}

func (builder *typecommonBuilder[T, S]) Primary() *S {
	builder.pairs = append(builder.pairs, pair{"primary", true})
	return builder.self()
}

func (builder *typecommonBuilder[T, S]) Nullable() *S {
	builder.pairs = append(builder.pairs, pair{"nullable", true})
	return builder.self()
}

func (builder *typecommonBuilder[T, S]) Default(dv T) *S {
	builder.pairs = append(builder.pairs, pair{"default", dv})
	return builder.self()
}

func (builder *typecommonBuilder[T, S]) DefaultExpr(expr string) *S {
	builder.pairs = append(builder.pairs, pair{"defaultexpr", expr})
	return builder.self()
}

func (builder *typecommonBuilder[T, S]) CheckExpr(expr string) *S {
	builder.pairs = append(builder.pairs, pair{"check", expr})
	return builder.self()
}

func (builder *typecommonBuilder[T, S]) Name(name string) *S {
	builder.pairs = append(builder.pairs, pair{"name", name})
	return builder.self()
}

func (builder *typecommonBuilder[T, S]) Comment(comment string) *S {
	builder.pairs = append(builder.pairs, pair{"comment", comment})
	return builder.self()
}

func (builder *typecommonBuilder[T, S]) Build() sqlx.SqlField {
	ins := &sqlx.FieldMetainfo{}

	for _, pair := range builder.pairs {
		switch pair.key {
		case "name":
			{
				ins.Name = pair.val.(string)
				break
			}
		case "unqiue":
			{
				ins.DdlOptions.Unique = true
				break
			}
		case "nullable":
			{
				ins.DdlOptions.Nullable = true
				break
			}
		case "default":
			{
				ins.DdlOptions.Default.Valid = true
				dv := reflect.ValueOf(pair.val)
				if dv.Kind() == reflect.String {
					// todo qoute
					ins.DdlOptions.Default.String = fmt.Sprintf(`'%s'`, pair.val.(string))
				} else {
					ins.DdlOptions.Default.String = fmt.Sprintf("%s", pair.val)
				}
				break
			}
		case "defaultexpr":
			{
				ins.DdlOptions.Default.Valid = true
				ins.DdlOptions.Default.String = pair.val.(string)
				break
			}
		case "autoincr":
			{
				ins.DdlOptions.AutoIncr = true
				break
			}
		case "primary":
			{
				ins.DdlOptions.PrimaryKey = true
				break
			}
		case "check":
			{
				ins.DdlOptions.Check = pair.val.(string)
				break
			}
		case "comment":
			{
				ins.DdlOptions.Comment = pair.val.(string)
				break
			}
		case "sqltype":
			{
				st := pair.val.(sqltype)
				ins.DdlOptions.SqlType = st.kind
				ins.DdlOptions.SqlTypeArgs = st.args
				break
			}
		}
	}
	return sqlx.SqlField{
		Ptr:      builder.ptr,
		Metainfo: ins,
	}
}
