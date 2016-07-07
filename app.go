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

func Auth(r *http.Request) (int, error) {
	token := r.Header.Get("Authorization")
	if token == "" {
		fmt.Printf("null token")
	}
	pt, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(encrypt_key), nil
	})
	if err != nil || pt.Valid != true {
		fmt.Printf("auth error")
	}
	fmt.Printf("%+v", pt.Claims)
	return 1, nil

}

func CreateToken(i int, b bool) (string, error) {
	jwtoken.Claims["user_id"] = i

	t, err := jwtoken.SignedString([]byte(encrypt_key))
	if err != nil {
		fmt.Println(err)
	}
	return t, nil
}

func (g GetHock) GetHandler(w http.ResponseWriter, r *http.Request) {
	_, err := Auth(r)
	if err != nil {
	}
	payload := mux.Vars(r)
	ok := g.validater(payload)
	if ok != true {
	}
	g.handler(w, r, payload)
}

type PostHock struct {
	handler func(http.ResponseWriter, *http.Request, map[string]string)
}

func (p PostHock) PostHandler(w http.ResponseWriter, r *http.Request) {
	_, err := Auth(r)
	if err != nil {
	}
	var payload map[string]string
	p.handler(w, r, payload)
}

func ValidPostPayload(p map[string]string) bool { return true }

func CreateTokenHandler(w http.ResponseWriter, r *http.Request) {
	id, err := model.ExistUser(r)
	if err != nil {
		return
	}
	t, err := CreateToken(id, false)
	if err != nil {
		return
	}
	fmt.Println(t)
	w.Header().Set("Authorization", t)
	return

}
func Test(w http.ResponseWriter, r *http.Request) {
	t, err := CreateToken(1, false)
	if err != nil {

	}
	fmt.Println(t)
	w.Header().Set("Authorization", t)
	fmt.Printf("\n\nw : %v\n\n\n", w)
	r.Header.Set("Authorization", t)
	fmt.Printf("\n\nr : %v\n\n\n", r)
	fmt.Printf(r.Header.Get("Authorization"))
	Auth(r)
	return

}
func DeleteTokenHandler(w http.ResponseWriter, r *http.Request) {
	id, err := Auth(r)
	if err != nil {
		return
	}
	t, err := CreateToken(id, true)
	if err != nil {
		return
	}
	w.Header().Set("Authorization", t)
	return
}

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
