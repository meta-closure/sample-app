package app

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"

	_ "github.com/go-sql-driver/mysql"
)

// to use JWT
type UserClaim struct {
	UserId int
	jwt.StandardClaims
}

func (l Login) Auth() (int, error) {
	pt, err := jwt.Parse(l.Token, func(token *jwt.Token) (interface{}, error) {

		// convert token parameter to map[string]interface{}
		tk, ok := token.Claims.(jwt.MapClaims)
		if ok != true {
			return nil, errors.Wrap(ErrInvalid, "Type convert to mapclaim error")
		}

		// get user_id from JWT
		id, ok := tk["UserId"].(float64)
		if ok != true {
			return nil, errors.Wrap(ErrTypeInvalid, "Type convert to user_id to float64")
		}

		// get salt from user_id to token validation check
		s := NewSalt(int(id), "")
		err := s.SelectById()
		if err != nil {
			return nil, errors.Wrapf(err, "User salt not exist: user_id: %v", int(id))
		}
		return []byte(s.Salt.String), nil
	})

	if err != nil {
		return 0, err
	}

	// can parse token and invalid
	if pt.Valid != true {
		return 0, ErrInvalidToken
	}

	tk, ok := pt.Claims.(jwt.MapClaims)
	if ok != true {
		return 0, errors.Wrap(ErrInvalid, "Type convert to mapclaim error")
	}

	id, ok := tk["UserId"].(float64)
	if ok != true {
		return 0, ErrInvalid
	}
	return int(id), nil
}

func (l *Login) Create(b []byte) error {
	u := &User{}
	err := u.FromJSON(b)
	if err != nil {
		return errors.Wrap(err, "Invalid JSON")
	}

	// user data exist check
	err = u.Get()
	if err != nil {
		return errors.Wrap(err, "User not exist")
	}

	// to get salt that generated for each user
	i := int(u.Id.Int64)
	s := NewSalt(i, "")
	err = s.SelectById()
	if err != nil {
		return errors.Wrapf(err, "User salt not exist: user_id: %v", i)
	}

	claim := UserClaim{
		i,
		jwt.StandardClaims{
			ExpiresAt: jwt.TimeFunc().AddDate(0, 0, 1).Unix(),
		},
	}

	jwtoken := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	t, err := jwtoken.SignedString([]byte(s.Salt.String))
	if err != nil {
		return errors.Wrapf(err, "JWT create error")
	}

	l.Token = t
	return nil
}
