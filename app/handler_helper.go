package app

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
	"golang.org/x/crypto/scrypt"
)

type Login struct {
	Token string `json:"token"`
}

type Message struct {
	Error string `json:"error"`
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

func NewLogin(t string) *Login {
	return &Login{Token: t}
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
