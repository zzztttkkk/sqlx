package sqlx

import "context"

type TableMetainfo struct {
	Name    string
	Fields  []*FieldMetainfo
	Indexes []*IndexMetainfo
	Options map[string]any
}

type ISqlDialectImpl interface {
	Migration(ctx context.Context, table *TableMetainfo) error
}
