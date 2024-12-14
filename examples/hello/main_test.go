package main

import (
	"context"
	"database/sql"
	"testing"

	jmsqlx "github.com/jmoiron/sqlx"
	"github.com/zzztttkkk/sqlx"
)

type Sum struct {
	Val1 int64 `db:"val1"`
	Val2 int64 `db:"val2"`
	Val3 int64 `db:"val3"`
	Val4 int64 `db:"val4"`
}

var (
	addsql = "select 1 + 23 as val1, 34 + 45 as val2, 45 + 56 as val3, 444 + 454 as val4"
)

func BenchmarkThisLib(b *testing.B) {
	db, _ := sql.Open("sqlite3", ":memory:")
	defer db.Close()
	stmt := sqlx.SelectStmt[struct{}](addsql, func() *Sum { return &Sum{} }).MustPrepare(context.Background(), db)
	for i := 0; i < b.N; i++ {
		stmt.MustQueryOne(context.Background(), nil)
	}
}

func BenchmarkJmSqlx(b *testing.B) {
	db, _ := jmsqlx.Connect("sqlite3", ":memory:")
	defer db.Close()

	var sum Sum

	for i := 0; i < b.N; i++ {
		db.Get(&sum, addsql)
	}
}
