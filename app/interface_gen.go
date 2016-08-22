package app

import (
	"encoding/json"

	"github.com/gocraft/dbr"
)

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

func (m *Posts) String() string {
	var s string
	return s
}

func (m *Post) String() string {
	var s string
	return s
}

func (m *User) String() string {
	var s string
	return s
}
func (m *Salt) String() string {
	var s string
	return s
}
