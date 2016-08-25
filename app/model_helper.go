package app

import "time"

func (p *Post) Insert() error {
	_, err := engine.Insert(p)
	if err != nil {
		return err
	}

	_, err = engine.Get(p)
	if err != nil {
		return err
	}
	return nil
}

func (u *User) Insert() error {
	_, err := engine.Insert(u)
	if err != nil {
		return err
	}
	_, err = engine.Get(u)
	if err != nil {
		return err
	}
	return nil
}

func (m *Salt) Insert() error {
	_, err := engine.Insert(m)
	if err != nil {
		return err
	}
	return nil
}

func NewSalt(i int, salt string) *Salt {
	s := &Salt{}
	s.Salt.Scan(salt)
	s.UserId.Scan(i)
	return s
}

func NewPost() *Post {
	post := &Post{}

	post.CreatedAt.Scan(time.Now().Unix())
	post.UpdatedAt.Scan(time.Now().Unix())

	return post
}

func NewUser() *User {
	user := &User{}

	user.CreatedAt.Scan(time.Time{}.Unix())
	user.UpdatedAt.Scan(time.Time{}.Unix())

	return user
}

func NewCircle() *Circle {
	circle := &Circle{}

	circle.CreatedAt.Scan(time.Time{}.Unix())
	circle.UpdatedAt.Scan(time.Time{}.Unix())

	return circle
}

func NewFollow() *FollowList {
	return &FollowList{}
}

func NewLikeList() *LikeList {
	return &LikeList{}
}

func NewJoinList() *JoinList {
	return &JoinList{}
}

func NewPostList() *PostList {
	return &PostList{}
}

func NewUserList() *UserList {
	return &UserList{}
}

func NewCircleList() *CircleList {
	return &CircleList{}
}

func (u *User) Select() error {
	ok, err := engine.Get(u)
	if err != nil {
		return err
	}
	if ok != true {
		return ErrEmpty
	}
	return nil
}
