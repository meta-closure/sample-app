package app

import (
	"fmt"
	"os"

	"github.com/go-xorm/xorm"
)

type TableMapper struct{}

func InitDb() *xorm.Engine {
	engine, err := xorm.NewEngine("mysql", "root:@/test")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	engine.SetMapper(TableMapper{})
	return engine
}

func (t TableMapper) Obj2Table(s string) string {
	switch s {
	case "Salt":
		return "salt"
	case "Post":
		return "posts"
	case "User":
		return "users"
	default:
		return "empty"
	}
}

func (t TableMapper) Table2Obj(s string) string {
	switch s {
	case "salt":
		return "Salt"
	case "posts":
		return "Post"
	case "users":
		return "User"
	default:
		return "empty"
	}
}
