package app

import (
	"fmt"

	"github.com/dgrijalva/jwt-go"
)

var (
	encrypt_key = "secret"
)

type UserClaim struct {
	UserId int
	jwt.StandardClaims
}

func Auth(token string) (int, error) {

	pt, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(encrypt_key), nil
	})
	if err != nil || pt.Valid != true {
		err = fmt.Errorf("authorization error")
		return 0, err
	}
	tk, _ := pt.Claims.(jwt.MapClaims)
	id, _ := tk["User_id"].(int)
	return id, nil
}

func CreateToken(i int) (string, error) {
	claim := UserClaim{
		i,
		jwt.StandardClaims{
			ExpiresAt: jwt.TimeFunc().AddDate(0, 0, 1).Unix(),
		},
	}
	jwtoken := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	t, err := jwtoken.SignedString([]byte(encrypt_key))
	if err != nil {
		fmt.Println(err)
	}
	return t, nil
}

func CreateExpiredToken(i int) (string, error) {
	claim := UserClaim{
		i,
		jwt.StandardClaims{
			ExpiresAt: 1,
		},
	}
	jwtoken := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	t, err := jwtoken.SignedString([]byte(encrypt_key))
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return t, nil
}
