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

type SumArgs struct {
	A int64 `db:"a"`
	B int64 `db:"b"`
}

type SumResult struct {
	Sum int64 `db:"val"`
}

// func (sr *SumResult) FieldPtrs() []any {
// 	return []any{&sr.Sum}
// }

func main() {
	db, _ := sql.Open("sqlite3", "file:demo.db")

	ctx := sqlx.WithDb(context.Background(), db)

	stmt := sqlx.SelectStmt[SumArgs](
		"select @a + @b + @a as val",
		func() *SumResult { return &SumResult{} },
	).
		MustPrepare(ctx, db).
		QueryLengthHint(1)

	var wg sync.WaitGroup
	var count = 1
	wg.Add(count)
	for i := 0; i < count; i++ {
		go func() {
			defer wg.Done()

			args := SumArgs{
				A: int64(rand.Int31n(100)),
				B: int64(rand.Int31n(100)),
			}
			result := stmt.MustQueryOne(ctx, &args)
			fmt.Println(args, result)
		}()
	}

	wg.Wait()
}
