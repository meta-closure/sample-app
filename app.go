package main

import (
	"fmt"
	"net/http"
	"time"

	"./model"
	"github.com/Sirupsen/logrus"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

var (
	log         = initLog()
	jwtoken     = jwt.New(jwt.SigningMethodHS256)
	encrypt_key = "secret"
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

func ParseJWT(s string) (int, error) {
	return 1, nil
}

func (g GetHock) GetHandler(w http.ResponseWriter, r *http.Request) {
	id, err := ParseJWT("test")
	if err != nil {
	}
	payload := mux.Vars(r)
	ok := g.validater(payload)
	if ok != true {
	}
	payload["user_id"] = string(id)
	g.handler(w, r, payload)
}

type PostHock struct {
	handler func(http.ResponseWriter, *http.Request, map[string]string)
}

func (p PostHock) PostHandler(w http.ResponseWriter, r *http.Request) {
	id, err := ParseJWT("test")
	if err != nil {
	}
	var payload map[string]string
	payload["user_id"] = string(id)
	p.handler(w, r, payload)
}

func ValidPostPayload(p map[string]string) bool { return true }

func CreateToken(w http.ResponseWriter, r *http.Request) {
	id, err := model.ExistUser(r)
	if err != nil {
		return
	}
	jwtoken.Claims["user_id"] = id
	t, err := jwtoken.SignedString([]byte(encrypt_key))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(t)
	w.Header().Set("Authorization", t)
	return

}

func DeleteToken(w http.ResponseWriter, r *http.Request) {}

func NoCheck(map[string]string) bool {
	return true
}

func (s *Server) SetupRoutes() {
	r := s.Router
	r.HandleFunc("/test", GetHock{handler: model.Tes, validater: NoCheck}.GetHandler)
	r.HandleFunc("/posts", GetHock{handler: model.GetPostsHandler, validater: NoCheck}.GetHandler).Methods("GET")
	r.HandleFunc("/posts/{post_id}", GetHock{handler: model.GetPostHandler, validater: NoCheck}.GetHandler).Methods("GET")
	r.HandleFunc("/posts", PostHock{handler: model.PostPostHandler}.PostHandler).Methods("POST")
	r.HandleFunc("/login", CreateToken)
	r.HandleFunc("/login", DeleteToken).Methods("DELETE")
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
