package main

import (
	"database/sql"

	"github.com/zzztttkkk/reflectx"
	"github.com/zzztttkkk/sqlx"
	"github.com/zzztttkkk/sqlx/sqltypes"
)

type CommonMix struct {
	Id        int64           `db:"id"`
	CreatedAt int64           `db:"created_at"`
	DeletedAt sql.Null[int64] `db:"deleted_at"`
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
	Name     string `db:"name"`
	Email    string `db:"email"`
	Password string `db:"password"`
}

func init() {
	mptr := reflectx.Ptr[User]()

	sqlx.Table[User]().
		Scheme().
		Tablename("account_user").
		Mixed(&CommonMix{}).
		Field(sqltypes.Varchar(&mptr.Name, 32).Unique().Build()).
		Field(sqltypes.Varchar(&mptr.Password, 155).Build())

	sqltypes.Varchar(&mptr.Email, 64).Build().AddToScheme(sqlx.Table[User]().Scheme())
	sqlx.Index[User]("email_index").Field(&mptr.Email, 0, nil)
}

type Post struct {
	CommonMix
	Uid     int64  `db:"uid"`
	Title   string `db:"title"`
	Content string `db:"content"`
}

func init() {
	mptr := reflectx.Ptr[Post]()
	sqlx.Table[Post]().Scheme().Tablename("post").Mixed(&CommonMix{})
	sqltypes.Int(&mptr.Uid).Build().AddToScheme(sqlx.Table[Post]().Scheme())
}
