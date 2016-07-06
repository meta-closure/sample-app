package model

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	_ "github.com/go-sql-driver/mysql"
)

const MYSQL_CONF = "root:@/test"

var (
	dbmap = initDb()
)

func initDb() *gorp.DbMap {
	db, err := sql.Open("mysql", MYSQL_CONF)
	if err != nil {
		os.Exit(-1)
	}
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"}}

	dbmap.AddTableWithName(Post{}, "posts")
	dbmap.AddTableWithName(User{}, "users")
	dbmap.TypeConverter = sampleTypeConverter{}
	//	dbmap.AddTableWithName(TestInt{}, "ok")
	return dbmap
}

func Tes(w http.ResponseWriter, r *http.Request, p map[string]string) {
	post := &Post{}
	err := dbmap.SelectOne(post, "select * from posts")
	fmt.Println(err)

	return
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

func Result(code int) []byte {
	ok := &Ok{Status: code}
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
	fmt.Printf("%v", user)
	if err != nil {
		fmt.Println(buf, err)
	}
	err = json.Unmarshal(buf, user)
	if err != nil {
		fmt.Println(err)
	}
	i, p := user.Id.Value().(int), user.CryptedPassword.Value().(string)
	err = dbmap.SelectOne(&User{}, "select * from users where id=? and crypted_password=?",
		i, p)
	if err != nil {
		fmt.Println(err)
	}
	return i, nil
}

func GetPostHandler(w http.ResponseWriter, r *http.Request, p map[string]string) {
	post := &Post{}
	err := dbmap.SelectOne(post, "select * from posts where id=?", p["post_id"])
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
	post := &Post{}
	posts := &Posts{}

	query, err := Query(r)
	if err != nil {
		fmt.Println(err)
	}
	if query["page"] == nil {
		query["page"] = []string{"1"}
	}
	j, err := strconv.Atoi(query["page"][0])
	if err != nil {
		fmt.Println(err)
	}
	st, ed := posts.Interval(j)
	s, err := dbmap.Select(post, "select * from posts where id between ? and ?",
		st, ed)
	if err != nil {
		fmt.Println(err)
	}
	for _, t := range s {
		posts.PostList = append(posts.PostList, t.(Post))
	}

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
	err = dbmap.Insert(user)
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
	err = dbmap.Insert(post)
	if err != nil {
		fmt.Printf("%+v", post)
		fmt.Println(err)
	}
	fmt.Fprintf(w, string(Result(200)))
	return
}
