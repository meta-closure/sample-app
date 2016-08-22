package app

import (
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/pkg/errors"
)

var (
	engine = InitDb()

	ErrInvalidPassword = errors.New("invalid password")
	ErrTypeInvalid     = errors.New("type error")
	ErrInvalid         = errors.New("parameter invalid")
	ErrEmpty           = errors.New("record is empty")
)

func NewSalt(i int, salt string) *Salt {
	s := &Salt{}
	s.Salt.Scan(salt)
	s.UserId.Scan(i)
	return s
}

func NewPosts() *Posts {
	return &Posts{}
}

func NewPost(b []byte) (*Post, error) {
	now := time.Time{}.Unix()
	post := &Post{}
	err := post.FromJSON(b)
	if err != nil {
		return post, err
	}
	post.CreatedAt.Scan(now)
	post.UpdatedAt.Scan(now)
	return post, nil
}

func NewUser(b []byte) (*User, error) {
	now := time.Time{}.Unix()
	user := &User{}
	err := user.FromJSON(b)
	if err != nil {
		return user, err
	}
	user.CreatedAt.Scan(now)
	user.UpdatedAt.Scan(now)

	return user, nil
}

func (s *Salt) SelectById() error {
	ok, err := engine.Where("user_id=?", int(s.UserId.Int64)).Get(s)
	if err != nil {
		return err
	}
	if ok != true {
		return ErrEmpty
	}
	return nil
}

func (u *User) Select() error {
	ok, err := engine.Get(u)
	if err != nil {
		return err
	}
	if ok != true {
		return ErrEmpty
	}
	return nil
}

func (u *User) Get() error {
	if u.Password.Valid {
		return errors.Wrap(ErrEmpty, "Password")
	}

	if u.Id.Valid {
		return errors.Wrap(ErrEmpty, "id")
	}

	err := u.Select()
	if err != nil {
		return err
	}

	return nil
}

func (p *Post) SelectById(id string) error {
	_, err := engine.Where("id=?", id).Get(p)
	if err != nil {
		return err
	}
	return nil
}

func (p *Posts) SelectByPage(j int) error {
	p.Page.Scan(j)
	pl := &[]Post{}
	st, ed := 1, 10
	err := engine.Where("id between ? and ?", st, ed).Find(pl)
	if err != nil {
		return err
	}
	p.PostList = *pl
	return nil
}

func (p *Post) Insert() error {
	_, err := engine.Insert(p)
	if err != nil {
		return err
	}
	return nil
}

func (u *User) Insert() error {
	_, err := engine.Insert(u)
	if err != nil {
		return err
	}
	_, err = engine.Get(u)
	if err != nil {
		return err
	}
	return nil
}

func (m *Salt) Insert() error {
	fmt.Println(m)
	_, err := engine.Insert(m)
	if err != nil {
		return err
	}
	return nil
}

func (p *Post) Update() error {
	return nil
}

func (u *User) Update() error {
	return nil
}

func (p *Post) Delete() error {
	return nil
}

func (u *User) Delete() error {
	return nil
}
