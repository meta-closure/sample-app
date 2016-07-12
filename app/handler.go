package app

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

func GetPostHandler(w http.ResponseWriter, r *http.Request, p map[string]string) {
	post := &Post{}
	err := post.SelectById(p["post_id"])
	if err != nil {
		Error(&w, 400, err)
		return
	}
	b, err := post.ToJSON()
	if err != nil {
		Error(&w, 400, err)
		return
	}
	Success(&w, b)
	return
}

func GetPostsHandler(w http.ResponseWriter, r *http.Request, p map[string]string) {
	query, err := Query(r)
	var j int
	if err != nil {
		Error(&w, 400, err)
		return
	}
	if query["page"] == nil {
		j = 1
	} else {
		j, err = strconv.Atoi(query["page"][0])
		if err != nil {
			Error(&w, 400, err)
			return
		}
	}
	posts := NewPosts()
	err = posts.SelectByPage(j)
	if err != nil {
		Error(&w, 400, err)
		return
	}
	b, err := posts.ToJSON()
	if err != nil {
		Error(&w, 400, err)
		return
	}
	Success(&w, b)
	return
}

func PostUserHandler(w http.ResponseWriter, r *http.Request) {
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		Error(&w, 400, err)
		return
	}
	user, err := NewUser(buf)
	if err != nil {
		Error(&w, 400, err)
		return
	}
	s, err := user.Pass2Hash()
	if err != nil {
		Error(&w, 400, err)
		return
	}
	err = user.Insert()
	if err != nil {
		Error(&w, 400, err)
		return
	}
	id, _ := user.Id.Value()
	tid, ok := id.(int64)
	fmt.Println(ok)
	salt := NewSalt(int(tid), s)
	salt.Insert()
	fmt.Println("ok4", "\n\n")
	b, err := user.ToJSON()
	Success(&w, b)
	return
}

func PostPostHandler(w http.ResponseWriter, r *http.Request, p map[string]string) {
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		Error(&w, 400, err)
		return
	}
	post, err := NewPost(buf)
	if err != nil {
		Error(&w, 400, err)
		return
	}
	ok := post.CheckValidUserId(p["auth_user_id"])
	if ok != nil {
		Error(&w, 400, ErrInvalid)
		return
	}
	err = post.Insert()
	if err != nil {
		Error(&w, 400, err)
		return
	}

	b, err := post.ToJSON()
	if err != nil {
		Error(&w, 400, err)
		return
	}
	Success(&w, b)
	return
}
