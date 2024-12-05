package sqlx

import (
	"fmt"
	"reflect"
	"testing"
)

type idMeta int

func (_ idMeta) SqlField() SqlField {
	return SqlField{
		Name:     "id",
		AutoIncr: true,
		SqlType:  "bigint",
		Primary:  true,
	}
}

type CommonMix[T any] struct {
	Id Field[int64, idMeta, T]
}

type nameMeta int

func (_ nameMeta) SqlField() SqlField {
	return SqlField{
		Name:    "name",
		Unique:  true,
		SqlType: "varchar(36)",
	}
}

type User struct {
	CommonMix[*User]
	Name Field[string, nameMeta, *User]
}

func TestRef(t *testing.T) {
	utype := reflect.TypeOf(User{})
	info := mktypeinfo(utype)
	fmt.Println(info)
	fmt.Println(fieldinfos)

	var user User
	fmt.Println(user.Id.SqlField(), typeofUser, typeofIdMeta)
}

var (
	typeofUser   = reflect.TypeOf(User{})
	typeofIdMeta = reflect.TypeOf(idMeta(0))
)

func BenchmarkReflectTypeOf(b *testing.B) {
	var user User
	for i := 0; i < b.N; i++ {
		user.Id.SqlField()
	}
}

func BenchmarkDirect(b *testing.B) {
	for i := 0; i < b.N; i++ {
		QuerySqlField(typeofUser, typeofIdMeta)
	}
}
