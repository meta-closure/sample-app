package app

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func ExistUser(r *http.Request) (int, error) {
	user := &User{}
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return 0, err
	}
	err = json.Unmarshal(buf, user)
	if err != nil {
		return 0, err
	}
	ok, err := engine.Get(user)
	if err != nil {
		return 0, err
	}
	if ok != true {
	}
	i, err := user.Id.Value()
	if err != nil {
		return 0, err
	}
	s, _ := i.(int64)
	return int(s), nil
}

func NewPost() *Post {
	now := time.Time{}.Unix()
	post := &Post{}
	post.CreatedAt.Scan(now)
	post.UpdatedAt.Scan(now)
	return post
}

func NewUser() *User {
	now := time.Time{}.Unix()
	user := &User{}
	user.CreatedAt.Scan(now)
	user.UpdatedAt.Scan(now)
	return user
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
