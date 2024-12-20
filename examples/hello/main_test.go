package main

import (
	"context"
	"database/sql"
	"testing"

	jmsqlx "github.com/jmoiron/sqlx"
	"github.com/zzztttkkk/sqlx"
)

type SumArgs struct {
	num1 int64
	num2 int64
	num3 int64
	num4 int64
}

type Sum struct {
	Val1 int64 `db:"val1"`
	Val2 int64 `db:"val2"`
	Val3 int64 `db:"val3"`
	Val4 int64 `db:"val4"`
}

var (
	AddNoArgsSql = "select 12 + 13 as val1, 14 + 15 as val2, 125 + 42 as val3, 86 + 2 as val4"
	AddArgsSql   = "select @num1 + @num2 as val1, @num1 + @num3 as val2, @num1 + @num4 as val3, @num2 + @num3 as val4"
)

func BenchmarkAutoArgs(b *testing.B) {
	db, _ := sql.Open("sqlite3", ":memory:")
	defer db.Close()
	stmt := sqlx.SelectStmt[SumArgs, Sum](AddArgsSql).MustPrepare(context.Background(), db)
	args := SumArgs{23, 34, 55, 567}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		stmt.MustQueryOne(context.Background(), &args)
	}
}

func BenchmarkThisLib(b *testing.B) {
	db, _ := sql.Open("sqlite3", ":memory:")
	defer db.Close()
	stmt := sqlx.SelectStmt[struct{}, Sum](AddNoArgsSql).MustPrepare(context.Background(), db)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		stmt.MustQueryOne(context.Background(), nil)
	}
}

func BenchmarkJmSqlx(b *testing.B) {
	db, _ := jmsqlx.Connect("sqlite3", ":memory:")
	defer db.Close()
	stmt, _ := db.PreparexContext(context.Background(), AddNoArgsSql)

	var sum Sum
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		stmt.Get(&sum)
	}
}
