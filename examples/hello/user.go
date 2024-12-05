package hello

import (
	"github.com/zzztttkkk/sqlx"
	"github.com/zzztttkkk/sqlx/sqltypes"
)

type idmeta int

//lint:ignore ST1006 /
func (_ idmeta) SqlField() sqlx.SqlField {
	return sqltypes.BigInt("id").Build()
}

type User struct {
	Id sqlx.Field[int64, idmeta, *User]
}

func init() {
	sqlx.RegisterTypeByValue(User{})
}
