package app

import (
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

var (
	ErrInvalidToken = errors.New("invalid auth token")
	ErrEmptyToken   = errors.New("auth token is empty")
)

type GetHock struct {
	handler   func(http.ResponseWriter, *http.Request, map[string]string)
	validater func(map[string]string) bool
}

type PostHock struct {
	handler   func(http.ResponseWriter, *http.Request, map[string]string)
	validater func(map[string]string) bool
}

type DeleteHock struct {
	handler   func(http.ResponseWriter, *http.Request, map[string]string)
	validater func(map[string]string) bool
}

type PutHock struct {
	handler   func(http.ResponseWriter, *http.Request, map[string]string)
	validater func(map[string]string) bool
}

func (g GetHock) GetHandler(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if token == "" {
		Failed(&w, r, 401, ErrEmptyToken)
		return
	}

	l := NewLogin(token)
	id, err := l.Auth()
	if err != nil {
		Failed(&w, r, 401, err)
		return
	}

	payload := mux.Vars(r)
	ok := g.validater(payload)
	if ok != true {

	}

	payload["auth_user_id"] = string(id)
	g.handler(w, r, payload)
	return
}

func (p PostHock) PostHandler(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if token == "" {
		Failed(&w, r, 401, ErrEmptyToken)
		return
	}

	l := NewLogin(token)
	id, err := l.Auth()
	if err != nil {
		Failed(&w, r, 401, ErrInvalidToken)
		return
	}

	payload := map[string]string{
		"auth_user_id": string(id),
	}

	p.handler(w, r, payload)
	return
}

func (p PutHock) PutHandler(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if token == "" {
		Failed(&w, r, 401, ErrEmptyToken)
		return
	}

	l := NewLogin(token)
	id, err := l.Auth()
	if err != nil {
		Failed(&w, r, 401, ErrInvalidToken)
		return
	}

	payload := map[string]string{
		"auth_user_id": string(id),
	}
	p.handler(w, r, payload)
	return
}

func (p DeleteHock) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if token == "" {
		Failed(&w, r, 401, ErrEmptyToken)
		return
	}

	l := NewLogin(token)
	id, err := l.Auth()
	if err != nil {
		Failed(&w, r, 401, ErrInvalidToken)
		return
	}

	payload := map[string]string{
		"auth_user_id": string(id),
	}
	p.handler(w, r, payload)
	return
}

func CreateTokenHandler(w http.ResponseWriter, r *http.Request) {
	l := NewLogin("")

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		l.Failed(&w, nil, err)
	}

	err = l.Create(b)
	if err != nil {
		l.Failed(&w, b, errors.Wrap(err, "Failed to create token"))
		return
	}

	l.Success(&w, b)
	return
}

func NoCheck(map[string]string) bool {
	return true
}

func (s *Server) SetupRoutes() {
	r := s.Router
	r.HandleFunc("/posts", GetHock{handler: GetPostsHandler, validater: NoCheck}.GetHandler).Methods("GET")
	r.HandleFunc("/posts/{post_id}", GetHock{handler: GetPostHandler, validater: NoCheck}.GetHandler).Methods("GET")
	r.HandleFunc("/posts", PostHock{handler: PostPostHandler, validater: NoCheck}.PostHandler).Methods("POST")
	r.HandleFunc("/login", CreateTokenHandler).Methods("POST")
	r.HandleFunc("/users", PostUserHandler).Methods("POST")
}

type Server struct {
	*mux.Router
}

func (s *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	origin := req.Header.Get("Origin")
	if origin != "" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Authorization")
	}
	if req.Method == "OPTION" {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	s.Router.ServeHTTP(w, req)
}

func Run(l string) error {
	r := New()
	return http.ListenAndServe(l, r)
}

func New() *Server {
	s := &Server{
		Router: mux.NewRouter(),
	}
	s.SetupRoutes()
	return s
}
