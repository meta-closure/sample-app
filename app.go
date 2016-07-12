package main

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"./app"
	"github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
)

var (
	log             = initLog()
	ErrInvalidToken = errors.New("invalid auth token")
	ErrEmptyToken   = errors.New("auth token is empty")
)

func initLog() *logrus.Logger {
	log := logrus.New()
	return log
}

func WithTime(l *logrus.Logger) *logrus.Entry {
	return l.WithField("TIME:", time.Now())
}

type GetHock struct {
	handler   func(http.ResponseWriter, *http.Request, map[string]string)
	validater func(map[string]string) bool
}

func (g GetHock) GetHandler(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if token == "" {
		app.Error(&w, 401, ErrEmptyToken)
		return
	}
	id, err := app.Auth(token)
	if err != nil {
		app.Error(&w, 401, ErrInvalidToken)
		return
	}
	payload := mux.Vars(r)
	ok := g.validater(payload)
	if ok != true {
	}
	payload["auth_user_id"] = string(id)
	g.handler(w, r, payload)
}

type PostHock struct {
	handler func(http.ResponseWriter, *http.Request, map[string]string)
}

func (p PostHock) PostHandler(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if token == "" {
		app.Error(&w, 401, ErrEmptyToken)
		return
	}
	id, err := app.Auth(token)
	if err != nil {
		app.Error(&w, 401, ErrInvalidToken)
		return
	}
	payload := map[string]string{
		"auth_user_id": string(id),
	}
	p.handler(w, r, payload)
}

func ValidPostPayload(p map[string]string) bool {
	return true
}

type DeleteHock struct {
	handler func(http.ResponseWriter, *http.Request, map[string]string)
}

type PutHock struct {
	handler func(http.ResponseWriter, *http.Request, map[string]string)
}

func CreateTokenHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.ExistUser(r)
	if err != nil {
		fmt.Println(err)
		app.Error(&w, 401, err)
		return
	}
	t, err := app.CreateToken(id)
	if err != nil {
		app.Error(&w, 401, err)
		return
	}
	app.LoginSuccess(&w, t)
}

func Test(w http.ResponseWriter, r *http.Request) {}

func NoCheck(map[string]string) bool {
	return true
}

func (s *Server) SetupRoutes() {
	r := s.Router
	r.HandleFunc("/test", Test)
	r.HandleFunc("/posts", GetHock{handler: app.GetPostsHandler, validater: NoCheck}.GetHandler).Methods("GET")
	r.HandleFunc("/posts/{post_id}", GetHock{handler: app.GetPostHandler, validater: NoCheck}.GetHandler).Methods("GET")
	r.HandleFunc("/posts", PostHock{handler: app.PostPostHandler}.PostHandler).Methods("POST")
	r.HandleFunc("/login", CreateTokenHandler).Methods("POST")
	r.HandleFunc("/users", app.PostUserHandler).Methods("POST")
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

func main() {
	WithTime(log).Info("okpk")
	Run(":8080")
}
