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

type Article struct {
	CommonMix[*Article]
}

var (
	user            User
	article         Article
	typeofUserId    = reflect.TypeOf(user.Id)
	typeofArticleId = reflect.TypeOf(article.Id)
)

func init() {
	RegisterTypeByValue(Article{}, User{})
}

func TestRef(t *testing.T) {
	fmt.Println(QuerySqlField(typeofArticleId) != nil)
	fmt.Println(QuerySqlField(typeofArticleId) == article.Id.SqlField())
	fmt.Println(QuerySqlField(typeofUserId) == user.Id.SqlField())
	fmt.Println(QuerySqlField(typeofUserId) != article.Id.SqlField())
}

func BenchmarkReflectTypeOf(b *testing.B) {
	var user User
	for i := 0; i < b.N; i++ {
		user.Id.SqlField()
	}
}

func BenchmarkDirect(b *testing.B) {
	for i := 0; i < b.N; i++ {
		QuerySqlField(typeofUserId)
	}
}
