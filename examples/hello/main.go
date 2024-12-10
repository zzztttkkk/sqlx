package main

import (
	"context"
	"database/sql"
	"fmt"
	"math/rand"
	"sync"

	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
	"github.com/zzztttkkk/sqlx"
)

type SumResult struct {
	Sum int64 `db:"val"`
}

func main() {
	db, _ := sql.Open("sqlite3", "file:demo.db")

	ctx := sqlx.WithDb(context.Background(), db)

	stmt := sqlx.SelectStmt(
		"select ? + ? as val",
		func() *SumResult { return &SumResult{} },
	).
		MustPrepare(ctx, db).
		QueryLengthHint(1)

	var wg sync.WaitGroup
	var count = 100
	wg.Add(count)
	for i := 0; i < count; i++ {
		go func() {
			defer wg.Done()

			a := rand.Int31n(1000)
			b := rand.Int31n(10000)

			v := stmt.MustQueryOne(ctx, a, b)
			fmt.Println(i, a, b, v, v.Sum == int64(a+b))
		}()
	}

	wg.Wait()
}
