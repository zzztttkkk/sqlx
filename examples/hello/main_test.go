package main

import (
	"context"
	"database/sql"
	"math/rand"
	"testing"

	jmsqlx "github.com/jmoiron/sqlx"
	"github.com/zzztttkkk/sqlx"
)

type SumArgs struct {
	Num1 int64 `db:"num1"`
	Num2 int64 `db:"num2"`
	Num3 int64 `db:"num3"`
	Num4 int64 `db:"num4"`
}

type Sum struct {
	Val1 int64 `db:"val1"`
	Val2 int64 `db:"val2"`
	Val3 int64 `db:"val3"`
	Val4 int64 `db:"val4"`
}

var (
	addsql = "select @num1 + @num2 as val1, @num1 + @num3 as val2, @num1 + @num4 as val3, @num2 + @num3 as val4"
)

func BenchmarkThisLib(b *testing.B) {
	db, _ := sql.Open("sqlite3", ":memory:")
	defer db.Close()
	stmt := sqlx.SelectStmt[SumArgs](addsql, func() *Sum { return &Sum{} }).MustPrepare(context.Background(), db)
	for i := 0; i < b.N; i++ {
		args := SumArgs{rand.Int63n(100), rand.Int63n(100), rand.Int63n(100), rand.Int63n(100)}
		stmt.MustQueryOne(context.Background(), &args)
	}
}

func BenchmarkJmSqlx(b *testing.B) {
	db, _ := jmsqlx.Connect("sqlite3", ":memory:")
	defer db.Close()

	stmt, _ := db.PreparexContext(context.Background(), addsql)

	var sum Sum

	for i := 0; i < b.N; i++ {
		args := []any{
			sql.Named("num1", rand.Int63n(100)),
			sql.Named("num2", rand.Int63n(100)),
			sql.Named("num3", rand.Int63n(100)),
			sql.Named("num4", rand.Int63n(100)),
		}
		stmt.Get(&sum, args...)
	}
}
