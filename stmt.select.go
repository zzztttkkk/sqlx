package sqlx

type op struct {
	left  any
	op    string
	right any
}

type dmlSelect[T ITable, A IArgs] struct {
	table [0]T
}

func (ds *dmlSelect[T, A]) Fields(fields ...ifaceField) *dmlSelect[T, A] {
	return ds
}
