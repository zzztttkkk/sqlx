package main

import (
	"reflect"

	"github.com/zzztttkkk/sqlx"
	"github.com/zzztttkkk/sqlx/sqltypes"
)

type User struct {
	Id    sqlx.Field[int64, idmeta, *User]
	Name  sqlx.Field[string, namemeta, *User]
	Email sqlx.Field[string, namemeta, *User]
}

func (_ User) TableMetaInfo() sqlx.TableMetaInfo {
	return sqltypes.Table("user").Build()
}

type idmeta int

func (_ idmeta) FieldMetaInfo() sqlx.FieldMetaInfo {
	return sqltypes.BigInt("id").Build()
}

type namemeta int

func (_ namemeta) FieldMetaInfo() sqlx.FieldMetaInfo {
	return sqltypes.Varchar("name", 30).Unique().Build()
}

func init() {
	tableofUser := sqlx.RegisterTable(User{})
	tableofUser.NewIndex("").AddField(User{}.Id, sqlx.OrderAsc)
}

func (user *User) ScanFields(stmttype reflect.Type) []sqlx.ISqlField {
	return []sqlx.ISqlField{
		&user.Email,
		&user.Id,
		&user.Name,
	}
}
