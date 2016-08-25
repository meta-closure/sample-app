package app

import (
	_ "github.com/go-sql-driver/mysql"

	"github.com/pkg/errors"
)

var (
	engine = InitDb()

	ErrInvalidPassword = errors.New("invalid password")
	ErrTypeInvalid     = errors.New("type error")
	ErrInvalid         = errors.New("parameter invalid")
	ErrEmpty           = errors.New("record is empty")
)

func (u *User) Get() error {
	if u.Password.Valid != true {
		return errors.Wrap(ErrEmpty, "Password")
	}

	if u.Id.Valid != true {
		return errors.Wrap(ErrEmpty, "id")
	}

	err := u.Select()
	if err != nil {
		return err
	}

	return nil
}

func (p *Post) SelectByPostId(pid int) error {
	ok, err := engine.Where("id=?", pid).Get(p)
	if err != nil {
		return err
	}

	if ok != true {
		return ErrInvalidPostId
	}
	return nil
}

func (s *Salt) SelectByUserId(uid int) error {
	ok, err := engine.Where("user_id=?", uid).Get(s)
	if err != nil {
		return err
	}

	if ok != true {
		return ErrInvalidUserId
	}

	return nil
}

func (u *User) SelectByUserId(uid int) error {
	ok, err := engine.Where("id=?", uid).Get(u)
	if err != nil {
		return err
	}

	if ok != true {
		return ErrInvalidUserId
	}
	return nil
}

func (c *Circle) SelectByCircleId(cid int) error {
	ok, err := engine.Where("id=?", cid).Get(c)
	if err != nil {
		return err
	}

	if ok != true {
		return ErrInvalidCircleId
	}
	return nil
}

func (u *UserList) Select(q Query) error {
	us := NewUserList()
	var err error
	if q.Order == "asc" {
		err = engine.Where("created_at<?", q.Time).Limit(q.Item).Asc("created_at").Find(us)
	} else {
		err = engine.Where("created_at<?", q.Time).Limit(q.Item).Desc("created_at").Find(us)
	}

	if err != nil {
		return err
	}
	return nil
}

func (c *CircleList) Select(q Query) error {
	var cs *[]Circle
	var err error
	if q.Order == "asc" {
		err = engine.Where("created_at<?", q.Time).Limit(q.Item).Asc("created_at").Find(cs)
	} else {
		err = engine.Where("created_at<?", q.Time).Limit(q.Item).Desc("created_at").Find(cs)
	}

	if err != nil {
		return err
	}

	var js *[]Join
	var us *[]User
	for i, circle := range *cs {
		err = engine.Where("circle_id=?", circle.Id.Int64).Limit(q.Item).Desc("created_at").Find(us)
		if err != nil {
			return err
		}
		ids := []int{}
		for _, j := range *js {
			ids = append(ids, int(j.UserId.Int64))
		}

		err = engine.In("id", ids).Find(us)
		if err != nil {
			return err
		}

		(*cs)[i].JoinedBy = *us
	}

	return nil
}

func (p *PostList) SelectByCircleId(cid int, q Query) error {
	ps := NewPostList()
	var err error
	if q.Order == "asc" {
		err = engine.Where("created_at<?", q.Time).And("circle_id=?", cid).Limit(q.Item).Join("INNER", "users", "users.id=posts.user_id").Asc("created_at").Find(ps)
	} else {
		err = engine.Where("created_at<?", q.Time).And("circle_id=?", cid).Limit(q.Item).Join("INNER", "users", "users.id=posts.user_id").Desc("created_at").Find(ps)
	}

	if err != nil {
		return err
	}
	return nil
}

func (p *PostList) SelectByUserId(uid int, q Query) error {
	ps := NewPostList()
	var err error
	if q.Order == "asc" {
		err = engine.Where("created_at<?", q.Time).Limit(q.Item).Asc("created_at").Find(ps)
	} else {
		err = engine.Where("created_at<?", q.Time).Limit(q.Item).Desc("created_at").Find(ps)
	}

	if err != nil {
		return err
	}
	return nil
}
