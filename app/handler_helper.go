package app

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"golang.org/x/crypto/scrypt"

	_ "github.com/go-sql-driver/mysql"
)

type Ok struct {
	Ok      bool   `json:"ok"`
	Message string `json:"message`
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

func Success(w *http.ResponseWriter, b []byte) {
	if b == nil {
		ok := Ok{Ok: true}
		b, _ = json.Marshal(ok)
	}
	(*w).WriteHeader(200)
	fmt.Fprintf(*w, string(b))
}

func Error(w *http.ResponseWriter, code int, err error) {
	ok := Ok{Ok: false,
		Message: fmt.Sprint(err),
	}
	b, _ := json.Marshal(ok)
	(*w).WriteHeader(code)
	fmt.Fprintf(*w, string(b))
}

func Pass2Hash(s string) string {
	c, _ := scrypt.Key([]byte(s), []byte("password"), 16384, 8, 1, 32)
	return hex.EncodeToString(c[:])
}

func Interval(l int) (int, int) {
	if l < 0 {
		l = 1
	}
	i := 10*(l) - 9
	j := 10*(l+1) - 10
	return i, j
}
