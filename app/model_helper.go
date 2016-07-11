package app

import (
	"errors"
	"fmt"
	"os"
	"regexp"

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

func (m Post) Valid(p map[string]string) error {

	title, err := m.Title.Value()
	if err != nil || title == nil {
		return ErrInvalid
	} else {
		ttitle, ok := title.(string)
		if ok != true {
			return ErrInvalid
		}
		if len(ttitle) > 255 {
			return errors.New("invalid title too long")
		}
	}

	body, err := m.Body.Value()
	if err != nil || body == nil {
		return ErrInvalid
	} else {
		tbody, ok := body.(string)
		if ok != true {
			return ErrInvalid
		}
		if len(tbody) > 20000 {
			return errors.New("invalid body too long")
		}
	}

	createdat, err := m.CreatedAt.Value()
	if err != nil || createdat == nil {
		return ErrInvalid
	} else {
		tcreatedat, ok := createdat.(string)
		if ok != true {
			return ErrInvalid
		}
		ok, err := regexp.MatchString("^\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}Z$", tcreatedat)
		if err != nil || ok != true {
			return errors.New("invalid created_at pattern")
		}
	}

	updatedat, err := m.UpdatedAt.Value()
	if err != nil || createdat == nil {
		return ErrInvalid
	} else {
		tupdatedat, ok := updatedat.(string)
		if ok != true {
			return ErrInvalid
		}
		ok, err := regexp.MatchString("^\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}Z$", tupdatedat)
		if err != nil || ok != true {
			return errors.New("invalid updated_at pattern")
		}
	}
	return nil
}
func (m User) Valid(p map[string]string) error {
	screenname, err := m.ScreenName.Value()
	if err != nil || screenname == nil {
		return ErrInvalid
	} else {
		tscreenname, ok := screenname.(string)
		if ok != true {
			return ErrInvalid
		}
		if len(tscreenname) > 255 {
			return errors.New("invalid tscreenname too long")
		}
	}
	password, err := m.Password.Value()
	if err != nil || password == nil {
		return ErrInvalid
	} else {
		tpassword, ok := password.(string)
		if ok != true {
			return ErrInvalid
		}
		if len(tpassword) < 8 {
			return errors.New("invalid password too short")
		}
		if len(tpassword) > 255 {
			return errors.New("invalid password too long")
		}
		ok, err := regexp.MatchString("^(?=.*?[a-z])(?=.*?[A-Z])(?=.*?\\d)[a-zA-Z\\d]*$", tpassword)
		if err != nil || ok != true {
			return errors.New("invalid password pattern")
		}
	}

	createdat, err := m.CreatedAt.Value()
	if err != nil || createdat == nil {
		return ErrInvalid
	} else {
		tcreatedat, ok := createdat.(string)
		if ok != true {
			return ErrInvalid
		}
		ok, err := regexp.MatchString("^\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}Z$", tcreatedat)
		if err != nil || ok != true {
			return errors.New("invalid created_at pattern")
		}
	}

	updatedat, err := m.UpdatedAt.Value()
	if err != nil || createdat == nil {
		return ErrInvalid
	} else {
		tupdatedat, ok := updatedat.(string)
		if ok != true {
			return ErrInvalid
		}
		ok, err := regexp.MatchString("^\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}Z$", tupdatedat)
		if err != nil || ok != true {
			return errors.New("invalid updated_at pattern")
		}
	}
	return nil
}
