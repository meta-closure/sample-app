package main

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"./auth"
	"./model"
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
		fmt.Printf("auth token is empty\n")
		return
	}
	id, err := auth.Auth(token)
	if err != nil {
		fmt.Println(err)
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
		fmt.Println(ErrEmptyToken)
		return
	}
	id, err := auth.Auth(token)
	if err != nil {
		fmt.Println(err)
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

func CreateTokenHandler(w http.ResponseWriter, r *http.Request) {
	id, err := model.ExistUser(r)
	if err != nil {
		fmt.Println(err)
		return
	}
	t, err := auth.CreateToken(id)
	if err != nil {
		fmt.Println(err)
		return
	}
	w.Header().Set("Authorization", t)
	return

}

func DeleteTokenHandler(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if token == "" {
		fmt.Printf("null token")
		return
	}
	id, err := auth.Auth(token)
	if err != nil {
		fmt.Println(err)
		return
	}
	t, err := auth.CreateExpiredToken(id)
	if err != nil {
		fmt.Println(err)
		return
	}
	w.Header().Set("Authorization", t)
	return
}

func Test(w http.ResponseWriter, r *http.Request) {}

func NoCheck(map[string]string) bool {
	return true
}

func (s *Server) SetupRoutes() {
	r := s.Router
	r.HandleFunc("/test", Test)
	r.HandleFunc("/posts", GetHock{handler: model.GetPostsHandler, validater: NoCheck}.GetHandler).Methods("GET")
	r.HandleFunc("/posts/{post_id}", GetHock{handler: model.GetPostHandler, validater: NoCheck}.GetHandler).Methods("GET")
	r.HandleFunc("/posts", PostHock{handler: model.PostPostHandler}.PostHandler).Methods("POST")
	r.HandleFunc("/login", CreateTokenHandler).Methods("POST")
	r.HandleFunc("/login", DeleteTokenHandler).Methods("DELETE")
	r.HandleFunc("/users", model.PostUserHandler).Methods("POST")
}

type Server struct {
	*mux.Router
}

func Run(l string) error {
	return http.ListenAndServe(l, New())
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
