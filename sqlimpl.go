package sqlx

import (
	"context"
	"database/sql"
)

type Driver interface {
	QueryAllTablesInDb(ctx context.Context, db *sql.DB) ([]*TableMetaInfo, error)
	GenerateMigrateSql(indb *TableMetaInfo, incode *TableMetaInfo) (string, error)
	GenerateTableDdlSql(table *TableMetaInfo) (string, error)
	GenerateIndexDdlSql(table *TableMetaInfo) (string, error)
	QueryParamPlaceholder(idx string) (string, error)
}

type IExector interface {
	Query(ctx context.Context, sql string, args ...any) (*sql.Rows, error)
	ExecSql(ctx context.Context, sql string, args ...any) (sql.Result, error)
}

func Exector() {
	
}
