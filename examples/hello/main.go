package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	jmsqlx "github.com/jmoiron/sqlx"
	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
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
	times  = 100000
)

func run_jmsqlx() {
	db, _ := jmsqlx.Connect("sqlite3", ":memory:")
	defer db.Close()

	var sum Sum

	begin := time.Now()

	for i := 0; i < times; i++ {
		db.Get(&sum, addsql)
	}

	fmt.Println("jmoiron", time.Since(begin).Nanoseconds())
}

func run_thislib() {
	db, _ := sql.Open("sqlite3", ":memory:")
	defer db.Close()

	stmt := sqlx.SelectStmt[struct{}](addsql, func() *Sum { return &Sum{} }).MustPrepare(context.Background(), db)

	begin := time.Now()

	for i := 0; i < times; i++ {
		stmt.MustQueryOne(context.Background(), nil)
	}

	fmt.Println("thislib", time.Since(begin).Nanoseconds())
}

func main() {
	run_jmsqlx()
	run_thislib()
}
