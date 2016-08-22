package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pkg/errors"
)

func WrapGetPostHandler(w http.ResponseWriter, r *http.Request) {
	var p map[string]string
	GetPostHandler(w, r, p)
}

func WrapPostPostHandler(w http.ResponseWriter, r *http.Request) {
	p := map[string]string{
		"auth_user_id": "3",
	}
	PostPostHandler(w, r, p)
}

func WrapUserPostHandler(w http.ResponseWriter, r *http.Request) {
	PostUserHandler(w, r)
}

func PostTestCase() io.Reader {
	p := &Post{}
	p.Title.Scan("test")
	p.Body.Scan("test body")
	p.UserId.Scan(3)

	b, err := json.Marshal(p)
	if err != nil {
		fmt.Println(err)
	}

	r := bytes.NewReader(b)
	return r
}

func UserTestCase() io.Reader {
	p := &User{}
	p.ScreenName.Scan("test_user")
	p.Password.Scan("test_pass")

	b, err := json.Marshal(p)
	if err != nil {
		fmt.Println(err)
	}

	r := bytes.NewReader(b)
	return r
}

func SetupMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/users", WrapUserPostHandler)
	mux.HandleFunc("/posts/0", WrapGetPostHandler)
	mux.HandleFunc("/posts", WrapPostPostHandler)

	return mux
}

func SetupTable() error {
	return nil
}

func TestHandler(t *testing.T) {
	mux := SetupMux()
	ts := httptest.NewServer(mux)
	defer ts.Close()

	// UserPostHandler test
	r := UserTestCase()
	res, err := http.Post(ts.URL+"/users", "application/json", r)
	if err != nil {
		t.Error(errors.Wrap(err, "UserPostHandler request"))
	} else {
		var b []byte
		_, err := res.Body.Read(b)
		if err != nil {
			t.Error(err)
		}
		fmt.Printf("%s", b)
	}

	// PostPostHandler test
	r = PostTestCase()
	res, err = http.Post(ts.URL+"/posts", "application/json", r)
	if err != nil {
		t.Error(errors.Wrap(err, "PostPostHandler request"))
	} else {
		var b []byte
		_, err := res.Body.Read(b)
		if err != nil {
			t.Error(err)
		}
		fmt.Printf("%s", b)
	}

	// GetPostHandler test
	res, err = http.Get(ts.URL + "/posts/0")
	if err != nil {
		t.Error(errors.Wrap(err, "GetPostHandler request"))
	} else {
		var b []byte
		_, err := res.Body.Read(b)
		if err != nil {
			t.Error(err)
		}
		fmt.Printf("%s", b)
	}
	return
}
