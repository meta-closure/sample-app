package model

import (
	"errors"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/gocraft/dbr"
)

var (
	ErrTypeInvalid = errors.New("type error")
	ErrInvalid     = errors.New("invalid")
	ErrEmpty       = errors.New("empty")
)

type TableMapper struct{}

type Post struct {
	Id        dbr.NullInt64  `xorm:"id",json:"post_id,omitempty"`
	CreatedAt dbr.NullInt64  `xorm:"created_at",json:"created_at,omitempty"`
	UpdatedAt dbr.NullInt64  `xorm:"updated_at",json:"updated_at,omitempty"`
	Title     dbr.NullString `xorm:"title",json:"title",json:"title"`
	Body      dbr.NullString `xorm:"body",json:"body",json:"body"`
	UserId    dbr.NullInt64  `xorm:"user_id",json:"user_id"`
}

type User struct {
	Id              dbr.NullInt64  `xorm:"id",json:"user_id,omitempty"`
	CreatedAt       dbr.NullInt64  `xorm:"created_at",json:"created_at,omitempty"`
	UpdatedAt       dbr.NullInt64  `xorm:"updated_at",json:"updated_at,omitempty"`
	ScreenName      dbr.NullString `xorm:"screen_name",json:"screen_name"`
	CryptedPassword dbr.NullString `xorm:"crypted_password"`
	Password        dbr.NullString `json:"password`
}

func (t TableMapper) Obj2Table(s string) string {
	switch s {
	case "Post":
		return "posts"
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
	}
	engine.SetMapper(TableMapper{})
	return engine
}

type Posts struct {
	PostList []Post        `json:"post_list"`
	Page     dbr.NullInt64 `json:"page"`
}

type Ok struct {
	Ok     bool `json:"ok"`
	Status int  `json:"status"`
}

func Interval(l int) (int, int) {
	if l < 0 {
		l = 1
	}
	i := 10*(l) - 9
	j := 10*(l+1) - 10
	return i, j
}

func NewPost() *Post {
	return &Post{}
}

func NewUser() *User {
	return &User{}
}
func (m Post) Valid(p map[string]string) bool {
	return true
}
func (m User) Valid(p map[string]string) bool {
	return true
}
