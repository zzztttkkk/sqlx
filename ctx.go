package sqlx

import (
	"context"
	"database/sql"
)

type ctxKeyType int

const (
	ctxKeyForDb = ctxKeyType(iota)
	ctxKeyForTx
)

func WithDb(ctx context.Context, db *sql.DB) context.Context {
	return context.WithValue(ctx, ctxKeyForDb, db)
}

func WithTx(ctx context.Context, tx *sql.Tx) context.Context {
	return context.WithValue(ctx, ctxKeyForTx, tx)
}
