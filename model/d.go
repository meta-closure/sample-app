package model

import (
	"errors"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/lestrrat/go-jsval"
	"gopkg.in/gorp.v1"
)

var (
	ErrTypeInvalid = errors.New("type error")
	ErrInvalid     = errors.New("invalid")
	ErrEmpty       = errors.New("empty")
)

type sampleTypeConverter struct{}

func (s sampleTypeConverter) ToDb(val interface{}) (interface{}, error) {
	switch val.(type) {
	case jsval.MaybeInt:
		m := val.(jsval.MaybeInt)
		if m.ValidFlag != true {
			return nil, nil
		}
		return m.Value(), nil
	case jsval.MaybeString:
		m := val.(jsval.MaybeString)
		if m.ValidFlag != true {
			return nil, nil
		}
		return m.Value(), nil
	case jsval.MaybeUint:
		m := val.(jsval.MaybeUint)
		if m.ValidFlag != true {
			return nil, nil
		}
		return m.Value(), nil
	case jsval.MaybeBool:
		m := val.(jsval.MaybeBool)
		if m.ValidFlag != true {
			return nil, nil
		}
		return m.Value(), nil
	default:
		return nil, ErrTypeInvalid
	}
}

func (s sampleTypeConverter) FromDb(target interface{}) (gorp.CustomScanner, bool) {
	
	switch target.(type) {
	case *jsval.MaybeInt:
		binder := func(holder, target interface{}) error {
			h, ok := holder.(*int)
			if ok != true {
				return ErrTypeInvalid
			}
			t, ok := target.(*jsval.MaybeInt)
			if ok != true {
				return ErrTypeInvalid
			}
			(*t).Set(*h)
			return nil
		}
		return gorp.CustomScanner{new(int), target, binder}, true
	case *jsval.MaybeString:
		binder := func(holder, target interface{}) error {
			h, ok := holder.(*string)
			if ok != true {
				return ErrTypeInvalid
			}
			t, ok := target.(*jsval.MaybeString)
			if ok != true {
				return ErrTypeInvalid
			}
			(*t).Set(*h)
			return nil
		}
		return gorp.CustomScanner{new(string), target, binder}, true

	case *jsval.MaybeBool:
		binder := func(holder, target interface{}) error {
			h, ok := holder.(*bool)
			if ok != true {
				return ErrTypeInvalid
			}
			t, ok := target.(*jsval.MaybeBool)
			if ok != true {
				return ErrTypeInvalid
			}
			(*t).Set(*h)
			return nil
		}
		return gorp.CustomScanner{new(bool), target, binder}, true
	case *jsval.MaybeUint:
		binder := func(holder, target interface{}) error {
			h, ok := holder.(*uint)
			if ok != true {
				return ErrTypeInvalid
			}
			t, ok := target.(*jsval.MaybeUint)
			if ok != true {
				return ErrTypeInvalid
			}
			(*t).Set(*h)
			return nil
		}
		return gorp.CustomScanner{new(uint), target, binder}, true
	}
	return gorp.CustomScanner{}, false
}

type Post struct {
	Id        jsval.MaybeInt    `db:"id",json:"post_id,omitempty"`
	CreatedAt jsval.MaybeInt    `db:"created_at",json:"created_at,omitempty"`
	UpdatedAt jsval.MaybeInt    `db:"updated_at",json:"updated_at,omitempty"`
	Title     jsval.MaybeString `db:"title",json:"title",json:"title"`
	Body      jsval.MaybeString `db:"body",json:"body",json:"body"`
	UserId    jsval.MaybeInt    `db:"user_id",json:"user_id"`
}

type User struct {
	Id              jsval.MaybeInt    `db:"id",json:"user_id,omitempty"`
	CreatedAt       jsval.MaybeInt    `db:"created_at",json:"created_at,omitempty"`
	UpdatedAt       jsval.MaybeInt    `db:"updated_at",json:"updated_at,omitempty"`
	ScreenName      jsval.MaybeString `db:"screen_name",json:"screen_name"`
	CryptedPassword jsval.MaybeString `db:"crypted_password"`
	Password        jsval.MaybeString `db:"-",json:"password`
}

type Posts struct {
	PostList []Post         `json:"post_list"`
	Page     jsval.MaybeInt `json:"page"`
}

type Ok struct {
	Ok     bool `json:"ok"`
	Status int  `json:"status"`
}

func (p *Posts) Interval(l int) (int, int) {
	i := 10*l + 1
	j := 11*l + 1
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
