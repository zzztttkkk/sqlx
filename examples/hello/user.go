package main

import (
	"database/sql"

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
		sqltypes.Int(&c.Id).Name("id").Primary().AutoIncr().Build(),
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
		Scheme().
		Tablename("account_user").
		Mixed(&CommonMix{}).
		Field(sqltypes.Varchar(&mptr.Name, 32).Unique().Build()).
		Field(sqltypes.Varchar(&mptr.Password, 155).Build())

	sqltypes.Varchar(&mptr.Email, 64).Name("email").Build().AddToScheme(sqlx.Table[User]().Scheme())

	sqlx.Index[User]("email_index").Field(&mptr.Email, 0, nil)
}

type Post struct {
	CommonMix
	Uid     int64
	Title   string
	Content string
}

func init() {
	mptr := sqlx.Mptr[Post]()
	sqlx.Table[Post]().Scheme().Tablename("post").Mixed(&CommonMix{})
	sqltypes.Int(&mptr.Uid).Name("uid").Build().AddToScheme(sqlx.Table[Post]().Scheme())
}
