package model

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

const MYSQL_CONF = "root:@/test"

var (
	engine = initDb()
)

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

func Result(code int) []byte {
	ok := &Ok{}
	if code == 200 {
		ok.Ok = true
	} else {
		ok.Ok = false
	}
	b, _ := json.Marshal(ok)
	return b
}

func ExistUser(r *http.Request) (int, error) {
	user := &User{}
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(buf, err)
	}
	err = json.Unmarshal(buf, user)
	if err != nil {
		fmt.Println(err)
	}
	ok, err := engine.Get(user)
	if err != nil {
		fmt.Println(err)
		fmt.Println(user)
	}
	if ok != true {
	}
	i, err := user.Id.Value()
	if err != nil {
		fmt.Println(err)
	}
	s, _ := i.(int64)
	return int(s), nil
}

func GetPostHandler(w http.ResponseWriter, r *http.Request, p map[string]string) {
	post := &Post{}
	_, err := engine.Where("id=?", p["post_id"]).Get(post)

	if err != nil {
		fmt.Println(err)
	}
	b, err := json.Marshal(post)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Fprintf(w, string(b))
	return
}

func GetPostsHandler(w http.ResponseWriter, r *http.Request, p map[string]string) {
	pl := &[]Post{}
	posts := &Posts{}
	query, err := Query(r)
	var j int
	if err != nil {
		fmt.Println(err)
	}
	if query["page"] == nil {
		j = 1
	} else {
		j, err = strconv.Atoi(query["page"][0])
		if err != nil {
			fmt.Println(err)
		}
	}
	st, ed := Interval(j)
	err = engine.Where("id between ? and ?", st, ed).Find(pl)
	fmt.Printf("%+v", pl)
	fmt.Println(st, ed)
	if err != nil {
		fmt.Println(err)
	}
	posts.PostList = *pl
	b, err := json.Marshal(posts)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Fprintf(w, string(b))
	return
}

func PostUserHandler(w http.ResponseWriter, r *http.Request) {
	user := NewUser()
	buf, err := ioutil.ReadAll(r.Body)
	fmt.Printf("%v", user)
	if err != nil {
		fmt.Println(buf, err)
	}
	err = json.Unmarshal(buf, user)
	if err != nil {
		fmt.Println(err)
	}
	_, err = engine.Insert(user)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Fprintf(w, string(Result(200)))
	return
}

func PostPostHandler(w http.ResponseWriter, r *http.Request, p map[string]string) {
	post := NewPost()
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("er")
		fmt.Println(err)
	}
	err = json.Unmarshal(buf, post)
	if err != nil {
		fmt.Printf("errrr")
		fmt.Println(err)
	}
	_, err = engine.Insert(post)
	if err != nil {
		fmt.Printf("%+v", post)
		fmt.Println(err)
	}
	fmt.Fprintf(w, string(Result(200)))
	return
}
