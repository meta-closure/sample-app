package app

import (
	"encoding/json"
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

type Salt struct {
	UserId dbr.NullInt64  `xorm:"user_id"`
	Salt   dbr.NullString `xorm:"salt"`
}

type Post struct {
	Id        dbr.NullInt64  `xorm:"id" json:"post_id,omitempty"`
	CreatedAt dbr.NullTime   `xorm:"created_at" json:"created_at,omitempty"`
	UpdatedAt dbr.NullTime   `xorm:"updated_at" json:"updated_at,omitempty"`
	Title     dbr.NullString `xorm:"title" json:"title"`
	Body      dbr.NullString `xorm:"body" json:"body"`
	UserId    dbr.NullInt64  `xorm:"user_id" json:"user_id"`
}

type User struct {
	Id              dbr.NullInt64  `xorm:"id" json:"user_id,omitempty"`
	CreatedAt       dbr.NullTime   `xorm:"created_at" json:"created_at,omitempty"`
	UpdatedAt       dbr.NullTime   `xorm:"updated_at" json:"updated_at,omitempty"`
	ScreenName      dbr.NullString `xorm:"screen_name" json:"screen_name"`
	CryptedPassword dbr.NullString `xorm:"crypted_password" json:"crypted_password"`
	Password        dbr.NullString `xorm:"-" json:"password"`
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

func initDb() *xorm.Engine {
	engine, err := xorm.NewEngine("mysql", "root:@/test")
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	engine.SetMapper(TableMapper{})
	return engine
}

func (m *User) ToJSON() ([]byte, error) {
	b, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return b, err
}

func (m *Post) ToJSON() ([]byte, error) {
	b, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return b, err
}

func (m *User) FromJSON(b []byte) error {
	err := json.Unmarshal(b, m)
	if err != nil {
		return err
	}
	return nil
}

func (m *Post) FromJSON(b []byte) error {
	err := json.Unmarshal(b, m)
	if err != nil {
		return err
	}
	return nil
}

func (m *Posts) FromJSON(b []byte) error {
	err := json.Unmarshal(b, m)
	if err != nil {
		return err
	}
	return nil
}

func (m *Posts) ToJSON() ([]byte, error) {
	b, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return b, err
}

func (u *User) Pass2Hash() (string, error) {
	if u.Password.Valid != true {
		return "", errors.New("password empty")
	}
	p := u.Password.String
	salt, err := CreateSalt()
	if err != nil {
		return "", err
	}
	u.CryptedPassword.Scan(Pass2Hash(p, salt))
	return salt, nil
}

func (u *User) ComparePassword(p string) error {
	if u.CryptedPassword.Valid != true {
		return ErrInvalid
	}
	hp := u.CryptedPassword.String
	if u.Id.Valid != true {
		return ErrInvalid
	}
	i := u.Id.Int64
	salt, err := SearchSaltById(int(i))
	if err != nil {
		return err
	}
	if Pass2Hash(p, salt) != hp {
		return ErrInvalidPassword
	}
	return nil
}
