package sqlx

import (
	"reflect"
)

type OrderKind int

const (
	OrderAsc = OrderKind(iota)
	OrderDesc
)

type IndexKey struct {
	Name  string
	Order OrderKind
}

type Index struct {
	Name    string
	Keys    []IndexKey
	Options map[string]any

	table *TableMetaInfo
}

func (idx *Index) AddField(field ifaceField, order OrderKind) *Index {
	sf := reflect.New(reflect.TypeOf(field)).Interface().(ISqlField).SqlField()
	idx.Keys = append(idx.Keys, IndexKey{
		Name:  sf.Name,
		Order: order,
	})
	return idx
}

func (idx *Index) Option(key string, val any) *Index {
	if idx.Options == nil {
		idx.Options = map[string]any{}
	}
	idx.Options[key] = val
	return idx
}

func (tmi *TableMetaInfo) NewIndex(name string) *Index {
	return &Index{
		table: tmi,
	}
}
