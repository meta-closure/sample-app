package app

import (
	"encoding/json"

	"github.com/gocraft/dbr"
)

type Salt struct {
	UserId dbr.NullInt64  `xorm:"user_id"`
	Salt   dbr.NullString `xorm:"salt"`
}

type Post struct {
	Id         dbr.NullInt64  `xorm:"id" json:"post_id,omitempty"`
	CreatedAt  dbr.NullTime   `xorm:"created_at" json:"created_at,omitempty"`
	UpdatedAt  dbr.NullTime   `xorm:"updated_at" json:"updated_at,omitempty"`
	Title      dbr.NullString `xorm:"title" json:"title"`
	Body       dbr.NullString `xorm:"body" json:"body"`
	UserId     dbr.NullInt64  `xorm:"user_id" json:"user_id"`
	Auther     User           `xorm:"-" json:"author"`
	LikedCount dbr.NullInt64  `xorm:"-" json:"liked_count"`
	LikedLink  []User         `xorm:"-" json:"liked_link,omitempty"`
	Medium     dbr.NullString `xorm:"medium" json:"medium,omitempty"`
}

type User struct {
	Id              dbr.NullInt64  `xorm:"id" json:"user_id,omitempty"`
	FollowedNum     dbr.NullInt64  `xorm:"followed_num" json:"followed_num"`
	CryptedPassword dbr.NullString `xorm:"crypted_password" json:"-"`
	Password        dbr.NullString `xorm:"-" json:"password,omitempty"`
	CreatedAt       dbr.NullTime   `xorm:"created_at" json:"created_at,omitempty"`
	UpdatedAt       dbr.NullTime   `xorm:"updated_at" json:"updated_at, omitempty"`
	Name            dbr.NullString `xorm:"name" json:"name,omitempty"`
	NameReading     dbr.NullString `xorm:"name_reading", json:"name_reading",omitempty"`
	Cover           dbr.NullString `xorm:"cover" json:"cover,omitempty"`
	Thumbnail       dbr.NullString `xorm:"thumbnail" json:"thumbnail,omitempty`
	Gender          dbr.NullString `xorm:"gender", json:"gender,omitempty"`
	Icon            dbr.NullString `xorm:"icon", json:"icon,omitempty"`
	Email           dbr.NullString `xorm:"email", json:"email,omitempty"`
	Profile         dbr.NullString `xorm:"profile", json:"profile,omitempty"`
}

type Circle struct {
	Id          dbr.NullInt64  `xorm:"id" json:"id,omitempty"`
	CreatedAt   dbr.NullTime   `xorm:"created_at" json:"created_at,omitempty"`
	UpdatedAt   dbr.NullTime   `xorm:"updated_at" json:"updated_at,omitempty"`
	Name        dbr.NullString `xorm:"name" json:"name,omitempty"`
	Description dbr.NullString `xorm:"description" json:"description,omitempty"`
	Cover       dbr.NullString `xorm:"cover" json:"cover,omitempty"`
	Thumbnail   dbr.NullString `xorm:"thumbnail" json:"thumbnail,omitempty"`
	JoinedBy    []User         `xorm:"-" json:"joined_by`
}

type Follow struct {
	FollowTo   dbr.NullInt64 `xorm:"follow_to"`
	FollowFrom dbr.NullInt64 `xorm:"follow_from"`
}

type Like struct {
	PostId dbr.NullInt64 `xorm:"post_id"`
	UserId dbr.NullInt64 `xorm:"user_id"`
}

type Join struct {
	UserId   dbr.NullInt64 `xrom:"user_id"`
	CircleId dbr.NullInt64 `xorm:"circle_id"`
}

type FollowList []Follow

type LikeList []Like

type JoinList []Join

type CircleList []Circle

type PostList []Post

type UserList []User

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

func (m *Circle) ToJSON() ([]byte, error) {
	b, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return b, err
}

func (m *CircleList) ToJSON() ([]byte, error) {
	b, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return b, err
}

func (m *UserList) ToJSON() ([]byte, error) {
	b, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return b, err
}

func (m *PostList) ToJSON() ([]byte, error) {
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

func (m *Circle) FromJSON(b []byte) error {
	err := json.Unmarshal(b, m)
	if err != nil {
		return err
	}
	return nil
}

func (m *CircleList) FromJSON(b []byte) error {
	err := json.Unmarshal(b, m)
	if err != nil {
		return err
	}
	return nil
}

func (m *PostList) FromJSON(b []byte) error {
	err := json.Unmarshal(b, m)
	if err != nil {
		return err
	}
	return nil
}

func (m *UserList) FromJSON(b []byte) error {
	err := json.Unmarshal(b, m)
	if err != nil {
		return err
	}
	return nil
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
