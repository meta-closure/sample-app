package app

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/scrypt"

	_ "github.com/go-sql-driver/mysql"
)

type Token struct {
	Token string `json:"token"`
}

type UserClaim struct {
	UserId int
	jwt.StandardClaims
}


func (u *User) Pass2Hash() (string, error) {
	if u.Password.Valid != true {
		return "", errors.New("password empty")
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

func (l *Login) Create(i int) error {
	claim := UserClaim{
		i,
		jwt.StandardClaims{
			ExpiresAt: jwt.TimeFunc().AddDate(0, 0, 1).Unix(),
		},
	}

	s, err := SearchSaltById(i)
	if err != nil {
		return err
	}

	jwtoken := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	t, err := jwtoken.SignedString([]byte(s))
	if err != nil {
		return err
	}

	l.Token = t
	return nil
}
