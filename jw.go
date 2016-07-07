package main

import (
	_ "encoding/json"
	"fmt"

	"github.com/dgrijalva/jwt-go"
)

type CustomClaim struct {
	UserId int
	jwt.StandardClaims
}

func main() {
	c := CustomClaim{
		1,
		jwt.StandardClaims{
			ExpiresAt: jwt.TimeFunc().Unix() - 10,
			Issuer:    "test",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)

	fmt.Printf("token: \n%+v\n", token)

	t, _ := token.SignedString([]byte("secret"))

	fmt.Printf("t: \n%+v\n", t)

	tk, _ := jwt.Parse(t, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})

	e, ok := tk.Claims.(jwt.MapClaims)
	fmt.Println(e.Valid(), ok)
	fmt.Printf("%+v", tk)

}
