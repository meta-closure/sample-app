package app

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/crypto/scrypt"
)

type Login struct {
	Token string `json:"token"`
}

type Message struct {
	Error string `json:"error"`
}

type Query struct {
	Item  int
	Time  int
	Order string
}

func NewQuery() Query {
	return Query{
		Item:  10,
		Time:  int(time.Now().Unix()),
		Order: "ASC",
	}
}

func ParseQuery(r *http.Request) (Query, error) {
	query := NewQuery()

	u, err := url.Parse(r.RequestURI)
	if err != nil {
		return query, err
	}

	q, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		return query, err
	}

	if q["item"] != nil {
		i, err := strconv.Atoi(q["item"][0])
		if err == nil {
			query.Item = i
		}
	}
	if q["time"] != nil {
		t, err := strconv.Atoi(q["time"][0])
		if err == nil {
			query.Time = t
		}
	}

	if q["sort"] != nil {
		if q["sort"][0] != "desc" {
			query.Order = "DESC"
		}
	}

	return query, nil
}

func NewLogin(t string) *Login {
	return &Login{Token: t}
}

func (l *Login) Success(w *http.ResponseWriter, r *http.Request, buf []byte) {
	b, _ := json.Marshal(l)
	Logger(w, r, 200, nil, buf, nil)
	(*w).WriteHeader(200)
	fmt.Fprintf(*w, string(b))
	return
}

func (l *Login) Failed(w *http.ResponseWriter, r *http.Request, buf []byte, err error) {
	Logger(w, r, 401, err, buf, nil)
	(*w).WriteHeader(401)
}

func Success(w *http.ResponseWriter, r *http.Request, buf, b []byte) {
	(*w).WriteHeader(200)
	Logger(w, r, 200, nil, buf, b)
	if b != nil {
		fmt.Fprintf(*w, string(b))
	}
	return
}

func Failed(w *http.ResponseWriter, r *http.Request, code int, err error) {
	b, _ := json.Marshal(Message{Error: err.Error()})
	(*w).WriteHeader(code)
	Logger(w, r, code, err, nil, nil)
	fmt.Fprintf(*w, string(b))
	return
}

func Logger(w *http.ResponseWriter, r *http.Request, code int, err error, req, res []byte) {
	var s string
	if err != nil {
		s = fmt.Sprintf("Host: %s%s\nMethod: %s\nHeader: %v\nStatusCode: %d\nError: %s\nRequest: %s\nResponse: %s\n",
			r.Host, r.RequestURI, r.Method, r.Header, code, err.Error(), req, res)
	} else {
		s = fmt.Sprintf("Host: %s%s\nMethod: %s\nHeader: %v\nStatusCode: %d\nRequest: %s\nResponse: %s\n",
			r.Host, r.RequestURI, r.Method, r.Header, code, req, res)
	}
	fmt.Println(s)
}

func (u *User) Pass2Hash() (string, error) {
	if u.Password.Valid != true {
		return "", errors.Wrap(ErrEmpty, "User password")
	}
	p := u.Password.String
	salt, err := CreateSalt()
	if err != nil {
		return "", err
	}
	u.CryptedPassword.Scan(Pass2Hash(p, salt))
	return salt, nil
}

func CreateSalt() (string, error) {
	b := make([]byte, 14)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	salt := base64.StdEncoding.EncodeToString(b)
	return salt, nil
}

func Pass2Hash(s, salt string) string {
	c, _ := scrypt.Key([]byte(s), []byte(salt), 16384, 8, 1, 32)
	return hex.EncodeToString(c[:])
}
