package app

import (
	"errors"

	"github.com/dgrijalva/jwt-go"
)

type UserClaim struct {
	UserId int
	jwt.StandardClaims
}

func Auth(token string) (int, error) {

	pt, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		tk, ok := token.Claims.(jwt.MapClaims)
		if ok != true {
			return nil, errors.New("invalid token")
		}
		id, ok := tk["UserId"].(float64)
		if ok != true {
			return nil, errors.New("invalid type error")
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
		return 0, errors.New("invalid token")
	}
	tk, _ := pt.Claims.(jwt.MapClaims)
	id, _ := tk["UserId"].(float64)
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
