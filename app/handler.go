package app

import (
	"io/ioutil"
	"net/http"

	"golang.org/x/net/context"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
)

var (
	ErrInvalidId   = errors.New("invalid user_id request")
	ErrInvalidAuth = errors.New("authorization failed")
)

func GetPostHandler(w http.ResponseWriter, r *http.Request, ctx context.Context) {
	post := &Post{}
	if ctx.Value("user_id") == nil {
		Failed(&w, r, 401, ErrInvalidAuth)
		return
	}

	id, ok := ctx.Value("user_id").(int)
	if ok != true {
		Failed(&w, r, 401, ErrInvalidAuth)
		return
	}

	err := post.SelectById(id)
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

func GetPostsHandler(w http.ResponseWriter, r *http.Request, ctx context.Context) {
	if ctx.Value("user_id") == nil {
		Failed(&w, r, 401, ErrInvalidAuth)
		return
	}

	q, err := ParseQuery(r)
	if err != nil {
		Failed(&w, r, 400, err)
		return
	}

	posts := NewPosts()
	err = posts.SelectByQuery(q)
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

func PostPostHandler(w http.ResponseWriter, r *http.Request, ctx context.Context) {
	if ctx.Value("user_id") == nil {
		Failed(&w, r, 401, ErrInvalidAuth)
		return
	}

	id, ok := ctx.Value("user_id").(int)
	if ok != true {
		Failed(&w, r, 401, ErrInvalidAuth)
		return
	}

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
	if id != int(post.UserId.Int64) {
		Failed(&w, r, 400, ErrInvalidId)
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
