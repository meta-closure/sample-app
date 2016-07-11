package app

import (
	"errors"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/gocraft/dbr"
)

var (
	ErrInvalidPassword = errors.New("invalid password")
	ErrTypeInvalid     = errors.New("type error")
	ErrInvalid         = errors.New("parameter invalid")
	ErrEmpty           = errors.New("record is empty")
)

type TableMapper struct{}

type Posts struct {
	PostList []Post        `json:"post_list"`
	Page     dbr.NullInt64 `json:"page"`
}

type Post struct {
	Id        dbr.NullInt64  `xorm:"id" json:"post_id,omitempty"`
	CreatedAt dbr.NullTime   `xorm:"created_at" json:"created_at omitempty"`
	UpdatedAt dbr.NullTime   `xorm:"updated_at" json:"updated_at omitempty"`
	Title     dbr.NullString `xorm:"title" json:"title" json:"title"`
	Body      dbr.NullString `xorm:"body" json:"body" json:"body"`
	UserId    dbr.NullInt64  `xorm:"user_id" json:"user_id"`
}

type User struct {
	Id              dbr.NullInt64  `xorm:"id" json:"user_id,omitempty"`
	CreatedAt       dbr.NullTime   `xorm:"created_at" json:"created_at omitempty"`
	UpdatedAt       dbr.NullTime   `xorm:"updated_at" json:"updated_at omitempty"`
	ScreenName      dbr.NullString `xorm:"screen_name" json:"screen_name"`
	CryptedPassword dbr.NullString `xorm:"crypted_password" json:"crypted_password"`
	Password        dbr.NullString `xorm:"-" json:"password"`
}

func (t TableMapper) Obj2Table(s string) string {
	switch s {
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
	case "posts":
		return "Post"
	case "users":
		return "User"
	default:
		return "empty"
	}
}

func initDb() *xorm.Engine {
	engine, err := xorm.NewEngine("mysql", "root:@/test")
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	engine.SetMapper(TableMapper{})
	return engine
}

func (u *User) Pass2Hash() error {
	p, err := u.Password.Value()
	if err != nil || p == nil {
		return ErrInvalid
	}
	ps, _ := p.(string)
	u.CryptedPassword.Scan(Pass2Hash(ps))
	fmt.Printf(Pass2Hash(ps))
	return nil
}

func (u *User) ComparePassword(p string) error {
	h, err := u.CryptedPassword.Value()
	if err != nil || h == nil {
		return ErrInvalid
	}
	hp, _ := h.(string)
	if Pass2Hash(p) != hp {
		return ErrInvalidPassword
	}
	return nil
}

func (m Post) Valid(p map[string]string) bool {
	return true
}
func (m User) Valid(p map[string]string) bool {
	return true
}
