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
		sqltypes.Int(&c.Id).Primary().AutoIncr().Build(),
		sqltypes.Int(&c.CreatedAt).DefaultExpr("unix_timestamp()").Build(),
		sqltypes.NullableInt(&c.DeletedAt).Nullable().Build(),
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
	sqlx.Table[User]().
		DDL().
		Tablename("name").
		Mixed(&CommonMix{}).
		Field(sqltypes.Varchar(&mptr.Name, 32).Unique().Build()).
		Field(sqltypes.Varchar(&mptr.Password, 155).Build())

	sqltypes.Varchar(&mptr.Email, 64).Build().AddToDdl(sqlx.Table[User]().DDL())

	sqlx.Index[User]("email_index").Field(unsafe.Pointer(&mptr.Email), 0, nil)
}
