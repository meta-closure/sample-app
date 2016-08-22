package app

import (
	"encoding/json"
)

type TableMapper struct{}

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
