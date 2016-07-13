package app

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"

	log "github.com/Sirupsen/logrus"
	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"golang.org/x/crypto/scrypt"

	_ "github.com/go-sql-driver/mysql"
)

type Ok struct {
	Ok      bool   `json:"ok"`
	Message string `json:"message`
}

type Token struct {
	Token string `json:"token"`
}

type UserClaim struct {
	UserId int
	jwt.StandardClaims
}

func Auth(token string) (int, error) {

	pt, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		tk, ok := token.Claims.(jwt.MapClaims)
		if ok != true {
			return nil, ErrInvalid
		}
		id, ok := tk["UserId"].(float64)
		if ok != true {
			return nil, ErrTypeInvalid
		}
		s, err := SearchSaltById(int(id))
		if err != nil {
			return nil, err
		}
		return []byte(s), nil
	})
	if err != nil {
		return 0, err
	}
	if pt.Valid != true {
		return 0, ErrInvalidToken
	}
	tk, ok := pt.Claims.(jwt.MapClaims)
	if ok != true {
		return 0, ErrInvalid
	}
	id, ok := tk["UserId"].(float64)
	if ok != true {
		return 0, ErrInvalid
	}
	return int(id), nil
}

func CreateToken(i int) (string, error) {
	claim := UserClaim{
		i,
		jwt.StandardClaims{
			ExpiresAt: jwt.TimeFunc().AddDate(0, 0, 1).Unix(),
		},
	}
	s, err := SearchSaltById(i)
	if err != nil {
		return "", err
	}
	jwtoken := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	t, err := jwtoken.SignedString([]byte(s))
	if err != nil {
		return "", err
	}
	return t, nil
}

func ExistUser(user *User) (int, error) {
	if user.Password.Valid {
		return 0, ErrInvalid
	}
	p := user.Password.String

	err := user.Select()
	if err != nil {
		return 0, err
	}

	err = user.ComparePassword(p)
	if err != nil {
		return 0, err
	}

	if user.Id.Valid != true {
		return 0, ErrInvalid
	}
	s := user.Id.Int64

	return int(s), nil
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

func LoginSuccess(w *http.ResponseWriter, r []byte, token string) {
	t := Token{Token: token}
	b, _ := json.Marshal(t)
	(*w).WriteHeader(200)
	log.Infof("Loggin Succeed, Request: %s", r)
	fmt.Fprintf(*w, string(b))
}

func LoginFailed(w *http.ResponseWriter, r []byte, err error) {
	(*w).WriteHeader(401)
	log.Infof("Loggin Failed, Request: %s, Error: %v", r, err)

}

func Success(w *http.ResponseWriter, r, b []byte) {
	(*w).WriteHeader(200)
	log.Infof("Request Succeed, Request: %s", r)
	if b != nil {
		fmt.Fprintf(*w, string(b))
	}
}

func Error(w *http.ResponseWriter, r []byte, code int, err error) {
	ok := Ok{Ok: false,
		Message: fmt.Sprint(err),
	}
	log.Infof("Request Failed, Request: %s, Error: %v", r, err)
	b, _ := json.Marshal(ok)
	(*w).WriteHeader(code)
	fmt.Fprintf(*w, string(b))
}

func Pass2Hash(s, salt string) string {
	c, _ := scrypt.Key([]byte(s), []byte(salt), 16384, 8, 1, 32)
	return hex.EncodeToString(c[:])
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
func SearchSaltById(i int) (string, error) {
	salt := &Salt{}
	err := salt.SelectById(i)
	if err != nil {
		errors.Wrapf(err, "Create &d' user id salt is fail", i)
		return "", err
	}
	s, _ := salt.Salt.Value()
	ts, _ := s.(string)
	return ts, nil
}

func Interval(l int) (int, int) {
	if l < 0 {
		l = 1
	}
	i := 10*(l) - 9
	j := 10*(l+1) - 10
	return i, j
}
