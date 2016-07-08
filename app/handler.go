package app

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

var (
	engine = initDb()
)

func GetPostHandler(w http.ResponseWriter, r *http.Request, p map[string]string) {
	post := &Post{}
	_, err := engine.Where("id=?", p["post_id"]).Get(post)

	if err != nil {
		Error(w, 400, err)
	}
	b, err := json.Marshal(post)
	if err != nil {
		Error(w, 400, err)
	}
	Success(w, b)
	return
}

func GetPostsHandler(w http.ResponseWriter, r *http.Request, p map[string]string) {
	pl := &[]Post{}
	posts := &Posts{}
	query, err := Query(r)
	var j int
	if err != nil {
		Error(w, 400, err)
	}
	if query["page"] == nil {
		j = 1
	} else {
		j, err = strconv.Atoi(query["page"][0])
		if err != nil {
			Error(w, 400, err)
		}
	}
	st, ed := Interval(j)
	err = engine.Where("id between ? and ?", st, ed).Find(pl)
	if err != nil {
		Error(w, 400, err)
	}
	posts.PostList = *pl
	b, err := json.Marshal(posts)
	if err != nil {
		Error(w, 400, err)
	}
	Success(w, b)
	return
}

func PostUserHandler(w http.ResponseWriter, r *http.Request) {
	user := NewUser()
	buf, err := ioutil.ReadAll(r.Body)
	fmt.Printf("%v", user)
	if err != nil {
		Error(w, 400, err)
	}
	err = json.Unmarshal(buf, user)
	if err != nil {
		Error(w, 400, err)
	}
	_, err = engine.Insert(user)
	if err != nil {
		Error(w, 400, err)
	}
	Success(w, nil)
	return
}

func PostPostHandler(w http.ResponseWriter, r *http.Request, p map[string]string) {
	post := NewPost()
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		Error(w, 400, err)
	}
	err = json.Unmarshal(buf, post)
	if err != nil {
		Error(w, 400, err)
	}
	_, err = engine.Insert(post)
	if err != nil {
		Error(w, 400, err)
	}
	Success(w, nil)
	return
}
