package sqlx

import (
	"context"
	"database/sql"
)

type InitOptions struct {
	AutoCreateTable        bool
	AutoCreateIndex        bool
	AllowDropUnusedTables  bool
	AllowDropUnusedIndexes bool
	AutoMigration          bool
}

func Init(ctx context.Context, db *sql.DB, driver Driver, opts *InitOptions) error {
	return nil
}
