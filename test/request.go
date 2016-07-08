package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/gocraft/dbr"
)

var (
	root = "http://localhost:8080"
)

type Auth struct {
	Token string
}

type User struct {
	Id              dbr.NullInt64  `xorm:"id",json:"user_id,omitempty"`
	CreatedAt       dbr.NullInt64  `xorm:"created_at" json:"created_at,omitempty"`
	UpdatedAt       dbr.NullInt64  `xorm:"updated_at" json:"updated_at,omitempty"`
	ScreenName      dbr.NullString `xorm:"screen_name" json:"screen_name"`
	CryptedPassword dbr.NullString `xorm:"crypted_password" json:"crypted_password"`
	Password        dbr.NullString `json:"password"`
}

func (a Auth) Request(method string, path string, p map[string]interface{}) {
	client := &http.Client{}
	var r io.Reader
	if method != "GET" {
		b, err := json.Marshal(p)
		if err != nil {
			fmt.Println(err)
		}
		r = bytes.NewReader(b)
	}
	url := root + path
	req, err := http.NewRequest(method, url, r)
	if err != nil {
		fmt.Println(err)
	}
	if a.Token != "" {
		req.Header.Add("Authorization", a.Token)
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(resp.Header.Get("Authorization"))
	fmt.Printf("%s", body)
}

func (a *Auth) Login(id int, pass string, create bool) error {
	client := &http.Client{}
	url := root + "/login"
	p := map[string]interface{}{
		"user_id":          id,
		"crypted_password": pass,
	}

	b, err := json.Marshal(p)
	if err != nil {
		fmt.Println(err)
	}
	r := bytes.NewReader(b)
	var method string
	if create == true {
		method = "POST"
	} else {
		method = "DELETE"
	}
	req, err := http.NewRequest(method, url, r)
	if err != nil {
		fmt.Println(err)
	}
	if a.Token != "" {
		req.Header.Add("Authorization", a.Token)
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	token := resp.Header.Get("Authorization")
	if token != "" {
		a.Token = token
		return nil
	} else {
		return err
	}
}

func main() {
	a := &Auth{}
	err := a.Login(1, "test", true)
	fmt.Println(a.Token)
	if err != nil {
		fmt.Println(err)
	}
	post := map[string]interface{}{
		"title": "fuck off",
		"body":  "off",
	}

	a.Request("GET", "/posts", post)
}
