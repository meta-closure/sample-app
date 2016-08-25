package app

import (
	"io/ioutil"
	"net/http"

	"golang.org/x/net/context"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
)

var (
	ErrInvalidPostId   = errors.New("invalid post_id request")
	ErrInvalidUserId   = errors.New("invalid user_id request")
	ErrInvalidCircleId = errors.New("invalid circle_id request")
	ErrInvalidAuth     = errors.New("authorization failed")
	ErrInvalidRequest  = errors.New("invalid request")
)

func GETCircleListHandler(w http.ResponseWriter, r *http.Request, ctx context.Context) {
	// if query parameter invalid, initialized value use
	query, _ := ParseQuery(r)

	circles := NewCircleList()
	err := circles.Select(query)
	if err != nil {
		Failed(&w, r, 400, ErrInvalidRequest)
		return
	}

	b, err := circles.ToJSON()
	if err != nil {
		Failed(&w, r, 400, ErrInvalidRequest)
		return
	}

	Success(&w, r, nil, b)
}

func GETCircleHandler(w http.ResponseWriter, r *http.Request, ctx context.Context) {
	param, err := ParseUrlParameter(r)
	if err != nil && param.CircleId == -1 {
		Failed(&w, r, 400, ErrInvalidCircleId)
		return
	}

	circle := NewCircle()
	err = circle.SelectByCircleId(param.CircleId)
	if err != nil {
		Failed(&w, r, 400, ErrInvalidCircleId)
		return
	}
	b, err := circle.ToJSON()
	if err != nil {
		Failed(&w, r, 400, ErrInvalidRequest)
		return
	}
	Success(&w, r, nil, b)
}

func GETCirclePostListHandler(w http.ResponseWriter, r *http.Request, ctx context.Context) {
	param, err := ParseUrlParameter(r)
	if err != nil && param.CircleId == -1 {
		Failed(&w, r, 400, ErrInvalidCircleId)
		return
	}

	query, _ := ParseQuery(r)

	posts := NewPostList()
	err = posts.SelectByCircleId(param.CircleId, query)
	if err != nil {
		Failed(&w, r, 400, ErrInvalidCircleId)
		return
	}

	b, err := posts.ToJSON()
	if err != nil {
		Failed(&w, r, 400, ErrInvalidRequest)
		return
	}
	Success(&w, r, nil, b)
}

// need auth
func GETMeHandler(w http.ResponseWriter, r *http.Request, ctx context.Context) {
	if ctx.Value("user_id") == nil {
		Failed(&w, r, 400, ErrInvalidAuth)
		return
	}

	uid, ok := ctx.Value("user_id").(int)
	if ok != true {
		Failed(&w, r, 400, ErrInvalidAuth)
		return
	}

	user := NewUser()
	err := user.SelectByUserId(uid)
	if err != nil {
		Failed(&w, r, 400, ErrInvalidUserId)
		return
	}

	b, err := user.ToJSON()
	if err != nil {
		Failed(&w, r, 400, ErrInvalidRequest)
		return
	}

	Success(&w, r, nil, b)
	return
}

func GETUserListHandler(w http.ResponseWriter, r *http.Request, ctx context.Context) {
	query, _ := ParseQuery(r)
	users := NewUserList()

	err := users.Select(query)
	if err != nil {
		Failed(&w, r, 400, ErrInvalidRequest)
		return
	}

	b, err := users.ToJSON()
	if err != nil {
		Failed(&w, r, 400, ErrInvalidRequest)
		return
	}

	Success(&w, r, nil, b)
	return
}
func GETUserHandler(w http.ResponseWriter, r *http.Request, ctx context.Context) {
	param, err := ParseUrlParameter(r)
	if err != nil && param.UserId == -1 {
		Failed(&w, r, 400, ErrInvalidUserId)
		return
	}

	user := NewUser()
	err = user.SelectByUserId(param.UserId)
	if err != nil && param.UserId == -1 {
		Failed(&w, r, 400, ErrInvalidUserId)
		return
	}

	b, err := user.ToJSON()
	if err != nil {
		Failed(&w, r, 400, ErrInvalidRequest)
		return
	}

	Success(&w, r, nil, b)
	return
}

func GETUserPostListHandler(w http.ResponseWriter, r *http.Request, ctx context.Context) {
	param, err := ParseUrlParameter(r)
	if err != nil && param.UserId == -1 {
		Failed(&w, r, 400, ErrInvalidUserId)
		return
	}

	query, _ := ParseQuery(r)
	posts := NewPostList()

	err = posts.SelectByUserId(param.UserId, query)
	if err != nil && param.UserId == -1 {
		Failed(&w, r, 400, ErrInvalidUserId)
		return
	}

	b, err := posts.ToJSON()
	if err != nil {
		Failed(&w, r, 400, ErrInvalidRequest)
		return
	}

	Success(&w, r, nil, b)
	return
}

func GETPostHandler(w http.ResponseWriter, r *http.Request, ctx context.Context) {
	param, err := ParseUrlParameter(r)
	if err != nil && param.PostId == -1 {
		Failed(&w, r, 400, ErrInvalidPostId)
		return
	}

	post := NewPost()

	err = post.SelectByPostId(param.PostId)
	if err != nil && param.UserId == -1 {
		Failed(&w, r, 400, ErrInvalidPostId)
		return
	}

	b, err := post.ToJSON()
	if err != nil {
		Failed(&w, r, 400, ErrInvalidRequest)
		return
	}

	Success(&w, r, nil, b)
	return
}

// old
func PostUserHandler(w http.ResponseWriter, r *http.Request) {
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		Failed(&w, r, 400, err)
		return
	}

	user := NewUser()
	err = user.FromJSON(buf)
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

// old
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

	post := NewPost()
	err = post.FromJSON(buf)
	if err != nil {
		Failed(&w, r, 400, err)
		return
	}

	// check request user id is user id
	if id != int(post.UserId.Int64) {
		Failed(&w, r, 400, ErrInvalidPostId)
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
