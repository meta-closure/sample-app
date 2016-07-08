package app

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

type Ok struct {
	Ok      bool  `json:"ok"`
	Message error `json:"message`
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

func Success(w http.ResponseWriter, b []byte) {
	if b == nil {
		ok := Ok{Ok: true}
		b, _ = json.Marshal(ok)
	}
	w.WriteHeader(200)
	fmt.Fprintf(w, string(b))
}

func Error(w http.ResponseWriter, code int, s error) {
	ok := Ok{Ok: false,
		Message: s,
	}
	b, _ := json.Marshal(ok)
	w.WriteHeader(code)
	fmt.Fprintf(w, string(b))
}

func Pass2Hash(s string) string {
	c, _ := bcrypt.GenerateFromPassword([]byte(s), 10)
	return hex.EncodeToString(c[:])
}
