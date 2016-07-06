package main

import (
	_ "encoding/json"
	"fmt"

	"github.com/dgrijalva/jwt-go"
)

func main() {
	token := jwt.New(jwt.SigningMethodHS256)

	token.Claims["foo"] = "bar"
	token.Claims["ee"] = "ok"
	fmt.Printf("token: \n%v\n", token)

	t, _ := token.SignedString([]byte("secret"))

	fmt.Printf("t: \n%+v\n", t)

	tk, _ := jwt.Parse(t, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	
}
