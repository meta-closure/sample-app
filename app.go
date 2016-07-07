package main

import (
	"fmt"
	"net/http"
	"time"

	"./auth"
	"./model"
	"github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
)

var (
	log = initLog()
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
		fmt.Printf("null token")
	}
	id, err := auth.Auth(token)
	if err != nil {
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
		fmt.Printf("null token")
	}
	id, err := auth.Auth(token)
	if err != nil {
	}
	var payload map[string]string
	payload["auth_user_id"] = string(id)
	p.handler(w, r, payload)
}

func ValidPostPayload(p map[string]string) bool { return true }

func CreateTokenHandler(w http.ResponseWriter, r *http.Request) {
	id, err := model.ExistUser(r)
	if err != nil {
		return
	}
	t, err := auth.CreateToken(id)
	if err != nil {
		return
	}
	fmt.Println(t)
	w.Header().Set("Authorization", t)
	return

}

func DeleteTokenHandler(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if token == "" {
		fmt.Printf("null token")
	}
	id, err := auth.Auth(token)
	if err != nil {
		return
	}
	t, err := auth.CreateExpiredToken(id)
	if err != nil {
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
	r.HandleFunc("/login", CreateTokenHandler)
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
