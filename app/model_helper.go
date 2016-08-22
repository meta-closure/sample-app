package app

import (
	"fmt"
	"os"

	"github.com/pkg/errors"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

var (
	ErrInvalidPassword = errors.New("invalid password")
	ErrTypeInvalid     = errors.New("type error")
	ErrInvalid         = errors.New("parameter invalid")
	ErrEmpty           = errors.New("record is empty")
)

func InitDb() *xorm.Engine {
	engine, err := xorm.NewEngine("mysql", "root:@/test")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	engine.SetMapper(TableMapper{})
	return engine
}
