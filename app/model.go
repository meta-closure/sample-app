package app

import (
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var (
	engine = initDb()
)

func (p Post) CheckValidUserId(s string) error {
	aud, _ := strconv.Atoi(s)
	u, err := p.UserId.Value()
	if err != nil || u == nil {
		return ErrInvalid
	}
	ud, ok := u.(int)
	if ok != true {
		return ErrInvalid
	}
	if aud != ud {
		return ErrInvalid
	}
	return nil
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
	st, ed := Interval(j)
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
	pass, err := u.Password.Value()
	if err != nil {
		return err
	}
	p := pass.(string)
	hash := Pass2Hash(p)
	u.CryptedPassword.Scan(hash)
	_, err = engine.Insert(u)
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
