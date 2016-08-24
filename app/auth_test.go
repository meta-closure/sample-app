package app

import (
	"encoding/json"
	"testing"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

type TestAuthCase struct {
	Case    map[string]interface{}
	Message string
	Pass    bool
}

func GetToken(i int, salt string) (string, error) {
	claim := UserClaim{
		i,
		jwt.StandardClaims{
			ExpiresAt: jwt.TimeFunc().AddDate(0, 0, 1).Unix(),
		},
	}

	jwtoken := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	t, err := jwtoken.SignedString([]byte(salt))
	if err != nil {
		return "", err
	}
	return t, nil
}

func GetExpiredToken(i int, salt string) (string, error) {
	claim := UserClaim{
		i,
		jwt.StandardClaims{
			ExpiresAt: jwt.TimeFunc().AddDate(0, 0, -1).Unix(),
		},
	}

	jwtoken := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	t, err := jwtoken.SignedString([]byte(salt))
	if err != nil {
		return "", err
	}
	return t, nil
}

func TestCreateToken(t *testing.T) {
	tests := []TestAuthCase{{
		Case: map[string]interface{}{
			"user_id":  45,
			"password": "test_pass",
		},
		Message: "JWT create pass",
		Pass:    true,
	}}

	for _, test := range tests {
		l := NewLogin("")
		u, err := json.Marshal(test.Case)
		if err != nil {
			t.Error(errors.Wrapf(err, "%s: encoding json", test.Message))
		}
		err = l.Create(u)
		// if test pass flag is true, token create and err should be empty
		if test.Pass == true && err != nil {
			t.Error(errors.Wrap(err, test.Message))
			continue
		}

		if test.Pass != true && err == nil {
			t.Error(test.Message)
			continue
		}
	}
}

func TestValidTokenAuth(t *testing.T) {
	// user_id. it is setted in DB before this test
	id := 45
	s := &Salt{}

	err := s.SelectById(id)
	if err != nil {
		t.Error(err)
	}

	token, err := GetToken(id, s.Salt.String)
	if err != nil {
		t.Error(err)
	}

	l := NewLogin(token)
	i, err := l.Auth()
	if err != nil {
		t.Error(err)
	}

	if i != id {
		t.Error("test should be pass")
	}
}

func TestExpiredTokenAuth(t *testing.T) {
	// user_id. it is setted in DB before this test
	id := 45
	s := &Salt{}

	err := s.SelectById(id)
	if err != nil {
		t.Error(err)
	}

	token, err := GetExpiredToken(id, s.Salt.String)
	if err != nil {
		t.Error(err)
	}

	l := NewLogin(token)
	_, err = l.Auth()
	if err == nil {
		t.Error(err)
	}
}

func TestInvalidSaltTokenAuth(t *testing.T) {
	// user_id. it is setted in DB before this test
	l := NewLogin("invalid token")
	_, err := l.Auth()
	if err == nil {
		t.Error(err)
	}
}
