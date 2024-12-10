package sqlx

import "context"

type TableMetainfo struct {
	Name    string
	Fields  []*FieldMetainfo
	Indexes []*IndexMetainfo
	Options map[string]any
}

type ISqlDialectImpl interface {
	Migrate(ctx context.Context, table *TableMetainfo) error
	QuoteValue(txt string) string
	QuoteName(name string) string
}
