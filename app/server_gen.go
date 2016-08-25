package app

import (
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

var (
	ErrInvalidToken = errors.New("invalid auth token")
	ErrEmptyToken   = errors.New("auth token is empty")
)

type AuthHock struct {
	handler   func(http.ResponseWriter, *http.Request, context.Context)
	validater func(map[string]string) bool
}

func (p AuthHock) AuthHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	token := r.Header.Get("Authorization")
	if token != "" {

		l := NewLogin(token)
		id, err := l.Auth()
		if err != nil {
			Failed(&w, r, 401, ErrInvalidToken)
			return
		}
		ctx = context.WithValue(ctx, "user_id", id)
	}

	p.handler(w, r, ctx)
	return
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {

	l := NewLogin("")
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		l.Failed(&w, r, b, err)
	}

	err = l.Create(b)
	if err != nil {
		l.Failed(&w, r, b, errors.Wrap(err, "Failed to create token"))
		return
	}

	l.Success(&w, r, b)
	return
}

func NoCheck(map[string]string) bool {
	return true
}

func (s *Server) SetupRoutes() {
	r := s.Router
	r.HandleFunc("/api/circles", AuthHock{handler: GETCircleListHandler}.AuthHandler).Methods("GET")
	r.HandleFunc("/api/circles/{circle_id}", AuthHock{handler: GETCircleHandler}.AuthHandler).Methods("GET")
	r.HandleFunc("/api/circles/{circle_id}/posts", AuthHock{handler: GETCirclePostListHandler}.AuthHandler).Methods("GET")
	r.HandleFunc("/api/posts/{post_id}", AuthHock{handler: GETPostHandler}.AuthHandler).Methods("GET")
	r.HandleFunc("/api/me", AuthHock{handler: GETMeHandler}.AuthHandler).Methods("GET")
	r.HandleFunc("/api/users", AuthHock{handler: GETUserListHandler}.AuthHandler).Methods("GET")
	r.HandleFunc("/api/users/{user_id}", AuthHock{handler: GETUserHandler}.AuthHandler).Methods("GET")
	r.HandleFunc("/api/users/{user_id}/posts", AuthHock{handler: GETUserPostListHandler}.AuthHandler).Methods("GET")
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
