package sqltypes

import "github.com/zzztttkkk/sqlx"

type tableMetaInfoBuilder struct {
	pairs []pair
}

func Table(name string) *tableMetaInfoBuilder {
	obj := &tableMetaInfoBuilder{}
	obj.pairs = append(obj.pairs, pair{"name", name})
	return obj
}

func (tmb *tableMetaInfoBuilder) Build() sqlx.TableMetaInfo {
	obj := sqlx.TableMetaInfo{}
	return obj
}
