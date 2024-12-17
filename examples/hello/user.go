package main

import (
	"database/sql"

	lion "github.com/zzztttkkk/reflectx"
	"github.com/zzztttkkk/sqlx"
	"github.com/zzztttkkk/sqlx/sqltypes"
)

type CommonMix struct {
	Id        int64           `db:"id"`
	CreatedAt int64           `db:"created_at"`
	DeletedAt sql.Null[int64] `db:"deleted_at"`
}

func init() {
	mptr := lion.Ptr[CommonMix]()
	schema := sqlx.Type[CommonMix]().Scheme()

	sqltypes.Int(&mptr.Id).Primary().Unique().AutoIncr().Build().AddToScheme(schema)
	sqltypes.Int(&mptr.CreatedAt).DefaultExpr("unix_timestamp()").Build().AddToScheme(schema)
	sqltypes.NullableInt(&mptr.DeletedAt).Nullable().Build().AddToScheme(schema)
}

type User struct {
	CommonMix
	Name     string `db:"name"`
	Email    string `db:"email"`
	Password string `db:"password"`
}

func init() {
	mptr := lion.Ptr[User]()

	sqlx.Type[User]().
		Scheme().
		Tablename("account_user").
		Field(sqltypes.Varchar(&mptr.Name, 32).Unique().Build()).
		Field(sqltypes.Varchar(&mptr.Password, 155).Build())

	sqltypes.Varchar(&mptr.Email, 64).Build().AddToScheme(sqlx.Type[User]().Scheme())
	sqlx.Index[User]("email_index").Field(&mptr.Email, 0, nil)
}

type Post struct {
	CommonMix
	Uid     int64  `db:"uid"`
	Title   string `db:"title"`
	Content string `db:"content"`
}

func init() {
	mptr := lion.Ptr[Post]()
	sqlx.Type[Post]().Scheme().Tablename("post")
	sqltypes.Int(&mptr.Uid).Build().AddToScheme(sqlx.Type[Post]().Scheme())
}
