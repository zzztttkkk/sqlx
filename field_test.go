package sqlx

import (
	"fmt"
	"reflect"
	"testing"
)

type idMeta int

func (_ idMeta) FieldMetaInfo() FieldMetaInfo {
	return FieldMetaInfo{
		Name: "id",
	}
}

type commonMix[T ITable] struct {
	Id Field[int64, idMeta, T]
}

type nameMeta int

func (_ nameMeta) FieldMetaInfo() FieldMetaInfo {
	return FieldMetaInfo{
		Name: "name",
	}
}

type User struct {
	commonMix[*User]
	Name Field[string, *nameMeta, *User]
}

func (_ User) TableMetaInfo() TableMetaInfo {
	return TableMetaInfo{}
}

type Article struct {
	commonMix[*Article]
}

func (_ Article) TableMetaInfo() TableMetaInfo {
	return TableMetaInfo{}
}

var (
	user            User
	article         Article
	typeofUserId    = reflect.TypeOf(user.Id)
	typeofArticleId = reflect.TypeOf(article.Id)
)

func init() {
	RegisterTable(User{})
	RegisterTable(Article{})
}

func TestRef(t *testing.T) {
	fmt.Println(QuerySqlField(typeofArticleId) != nil)
	fmt.Println(QuerySqlField(typeofArticleId) == article.Id.SqlField())
	fmt.Println(QuerySqlField(typeofUserId) == user.Id.SqlField())
	fmt.Println(QuerySqlField(typeofUserId) != article.Id.SqlField())
	fmt.Println(user.Name.SqlField())
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
