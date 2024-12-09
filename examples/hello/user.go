package main

import (
	"database/sql"
	"unsafe"

	"github.com/zzztttkkk/sqlx"
	"github.com/zzztttkkk/sqlx/sqltypes"
)

type CommonMix struct {
	Id        int64
	CreatedAt int64
	DeletedAt sql.Null[int64]
}

func (c *CommonMix) MixedFields() []sqlx.SqlField {
	return []sqlx.SqlField{
		sqltypes.BigInt(&c.Id, "id").Primary().AutoIncr().Build(),
		sqltypes.BigInt(&c.CreatedAt, "created_at").DefaultExpr("unix_timestamp()").Build(),
		sqltypes.BigInt(sqltypes.CastPtr[int64](unsafe.Pointer(&c.DeletedAt)), "deleted_at").Nullable().Build(),
	}
}

var _ sqlx.IMixed = (*CommonMix)(nil)

type User struct {
	CommonMix
	Name     string
	Email    string
	Password string
}

func init() {
	mptr := sqlx.Mptr[User]()
	sqlx.Table[User]().DDL().Tablename("name").
		Mixed(&CommonMix{}).
		Field(sqltypes.Varchar(&mptr.Name, "name", 32).Unique().Build()).
		Field(sqltypes.Varchar(&mptr.Email, "email", 64).Unique().Build()).
		Field(sqltypes.Varchar(&mptr.Password, "pwd", 155).Build()).
		Finish()

	sqltypes.Varchar(&mptr.Email, "email", 34).Build().AddToDdl(sqlx.Table[User]().DDL())
}
