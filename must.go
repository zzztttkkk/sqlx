package sqlx

func must[T any](v T, err error) T {
	if err == nil {
		return v
	}
	panic(err)
}
