package app

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"golang.org/x/net/context"

	"github.com/pkg/errors"
)

type TestCase struct {
	Case    map[string]interface{}
	Message string
	Url     string
	Pass    bool
}

type TestHandleWrapper struct {
	handler func(http.ResponseWriter, *http.Request, context.Context)
}

func (t TestHandleWrapper) WrapHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, "user_id", 3)
	t.handler(w, r, ctx)
}

func EncodeJSON(p map[string]interface{}) (io.Reader, error) {
	b, err := json.Marshal(p)
	if err != nil {
		return nil, errors.Wrap(err, "parameter encode to JSON")
	}

	r := bytes.NewReader(b)
	return r, nil
}

func SetupMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/users", PostUserHandler)
	mux.HandleFunc("/posts/1", TestHandleWrapper{handler: GetPostHandler}.WrapHandler)
	mux.HandleFunc("/posts", TestHandleWrapper{handler: PostPostHandler}.WrapHandler)

	return mux
}

func SetupTable() error {
	return nil
}

func TestHandlerStatus(t *testing.T) {
	mux := SetupMux()
	ts := httptest.NewServer(mux)
	defer ts.Close()

	tests := []TestCase{{
		Case: map[string]interface{}{
			"screen_name": "test_user",
			"password":    "test_pass",
		},
		Url:     "/users",
		Message: "UserPostHandler request Pass",
		Pass:    true,
	}, {
		Case: map[string]interface{}{
			"title":   "test title",
			"body":    "test body",
			"user_id": 3,
		},
		Url:     "/posts",
		Message: "PostPostHandler request Pass",
		Pass:    true,
	}, {
		Case: map[string]interface{}{
			"title":   "test title",
			"body":    "test body",
			"user_id": 5,
		},
		Url:     "/posts",
		Message: "PostPostHandler request Fail, invalid user_id",
		Pass:    false,
	}, {
		Case:    map[string]interface{}{},
		Url:     "/posts/1",
		Message: "GetPostHandler request Pass",
		Pass:    true,
	}}

	for _, test := range tests {

		r, err := EncodeJSON(test.Case)
		if err != nil {
			t.Error(errors.Wrap(err, test.Message))
		}
		res, err := http.Post(ts.URL+test.Url, "application/json", r)
		// request error
		if err != nil {
			t.Error(errors.Wrap(err, test.Url))
			continue
		}

		b, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Error(errors.Wrap(err, test.Url))
			continue
		}

		// if test pass flag true, then request should be success
		if test.Pass == true && res.StatusCode != 200 {
			t.Errorf("%s, Error: %s", test.Message, b)
			continue
		}

		// if test pass flag false, then request should not be pass
		if test.Pass != true && res.StatusCode == 200 {
			t.Errorf("%s, Error: %s", test.Message, b)
			continue
		}
	}
	return
}
