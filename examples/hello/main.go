package main

import (
	"fmt"

	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
	"github.com/zzztttkkk/sqlx"
)

func main() {
	tableinfo := sqlx.Type[User]().Scheme().Finish()
	fmt.Println(tableinfo)
}
