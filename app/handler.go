package app

import (
	"io/ioutil"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
)

func GetPostHandler(w http.ResponseWriter, r *http.Request, p map[string]string) {
	post := &Post{}
	err := post.SelectById(p["post_id"])
	if err != nil {
		Failed(&w, r, 400, err)
		return
	}
	b, err := post.ToJSON()
	if err != nil {
		Failed(&w, r, 400, err)
		return
	}
	Success(&w, r, nil, b)
	return
}

func GetPostsHandler(w http.ResponseWriter, r *http.Request, p map[string]string) {
	query, err := Query(r)
	var j int
	if err != nil {
		Failed(&w, r, 400, err)
		return
	}
	if query["page"] == nil {
		j = 1
	} else {
		j, err = strconv.Atoi(query["page"][0])
		if err != nil {
			Failed(&w, r, 400, err)
			return
		}
	}
	posts := NewPosts()
	err = posts.SelectByPage(j)
	if err != nil {
		Failed(&w, r, 400, err)
		return
	}
	b, err := posts.ToJSON()
	if err != nil {
		Failed(&w, r, 400, err)
		return
	}
	Success(&w, r, nil, b)
	return
}

func PostUserHandler(w http.ResponseWriter, r *http.Request) {
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		Failed(&w, r, 400, err)
		return
	}

	user, err := NewUser(buf)
	if err != nil {
		Failed(&w, r, 400, err)
		return
	}

	s, err := user.Pass2Hash()
	if err != nil {
		Failed(&w, r, 400, err)
		return
	}
	err = user.Insert()
	if err != nil {
		Failed(&w, r, 400, err)
		return
	}

	salt := NewSalt(int(user.Id.Int64), s)
	salt.Insert()
	b, err := user.ToJSON()
	Success(&w, r, buf, b)
	return
}

func PostPostHandler(w http.ResponseWriter, r *http.Request, p map[string]string) {
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		Failed(&w, r, 400, err)
		return
	}

	post, err := NewPost(buf)
	if err != nil {
		Failed(&w, r, 400, err)
		return
	}

	// check request user id is user id
	i, _ := strconv.Atoi(p["auth_user_id"])
	if i != int(post.UserId.Int64) {
		Failed(&w, r, 400, errors.New("User id not exist"))
		return
	}

	err = post.Insert()
	if err != nil {
		Failed(&w, r, 400, err)
		return
	}

	b, err := post.ToJSON()
	if err != nil {
		Failed(&w, r, 400, err)
		return
	}
	Success(&w, r, buf, b)
	return
}
