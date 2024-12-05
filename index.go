package sqlx

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
	Kind    string
	Keys    []IndexKey
	Options map[string]string
	Table   *SqlTable
}

func (index *Index) AddField(field *SqlField) {}
