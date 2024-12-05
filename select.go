package sqlx

type dmlSelect struct {
}

func (ds *dmlSelect) Fields(fields ...*SqlField) *dmlSelect {
	return ds
}
