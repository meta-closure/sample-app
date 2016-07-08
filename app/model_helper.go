package app

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

func ExistUser(r *http.Request) (int, error) {
	user := &User{}
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return 0, err
	}
	err = json.Unmarshal(buf, user)
	if err != nil {
		return 0, err
	}
	ok, err := engine.Get(user)
	if err != nil {
		return 0, err
	}
	if ok != true {
	}
	i, err := user.Id.Value()
	if err != nil {
		return 0, err
	}
	s, _ := i.(int64)
	return int(s), nil
}
