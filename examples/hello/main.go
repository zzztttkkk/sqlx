package main

import (
	"fmt"

	"github.com/zzztttkkk/sqlx"
)

func main() {
	fmt.Println(sqlx.Table[User]())
}
