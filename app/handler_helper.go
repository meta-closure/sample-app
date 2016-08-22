package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type Login struct {
	Token string
}

type Message struct {
	Error string
}

func (l Login) Valid() bool {
	return true
}

func Query(r *http.Request) (map[string][]string, error) {
	u, err := url.Parse(r.RequestURI)
	if err != nil {
		return nil, err
	}
	j, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		return nil, err
	}
	return j, nil
}

func LoginNew() *Login {
	return &Login{}
}

func (l *Login) Success(w *http.ResponseWriter, r []byte) {
	b, _ := json.Marshal(l)
	(*w).WriteHeader(200)
	fmt.Fprintf(*w, string(b))
	return
}

func (l *Login) Failed(w *http.ResponseWriter, r []byte, err error) {
	(*w).WriteHeader(401)
}

func Success(w *http.ResponseWriter, r, b []byte) {
	(*w).WriteHeader(200)
	if b != nil {
		fmt.Fprintf(*w, string(b))
	}
	return
}

func Failed(w *http.ResponseWriter, r *http.Request, code int, err error) {
	b, _ := json.Marshal(Message{Error: err.Error()})
	(*w).WriteHeader(code)
	fmt.Fprintf(*w, string(b))
	return
}
