package app

import (
	"errors"
	"fmt"
	"regexp"
)

func (m Post) Valid(p map[string]string) error {

	if m.Title.Valid != true {
		return ErrInvalid
	} else {
		title := m.Title.String
		if len(title) > 255 {
			return errors.New("invalid title too long")
		}
	}
	if m.Body.Valid != true {
		return ErrInvalid
	} else {
		body := m.Body.String
		if len(body) > 20000 {
			return errors.New("invalid body too long")
		}
	}

	if m.CreatedAt.Valid != true {
		return ErrInvalid
	} else {
		createdat := fmt.Sprint(m.CreatedAt.Time)
		ok, err := regexp.MatchString("^\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}Z$", createdat)
		if err != nil || ok != true {
			return errors.New("invalid created_at pattern")
		}
	}

	if m.UpdatedAt.Valid != true {
		return ErrInvalid
	} else {
		updatedat := fmt.Sprint(m.UpdatedAt.Time)
		ok, err := regexp.MatchString("^\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}Z$", updatedat)
		if err != nil || ok != true {
			return errors.New("invalid updated_at pattern")
		}
	}
	return nil
}

func (m User) Valid(p map[string]string) error {

	if m.Name.Valid != true {
		return ErrInvalid
	} else {
		name := m.Name.String
		if len(name) > 255 {
			return errors.New("invalid tscreenname too long")
		}
	}

	if m.Password.Valid != true {
		return ErrInvalid
	} else {
		password := m.Password.String
		if len(password) < 8 {
			return errors.New("invalid password too short")
		}
		if len(password) > 255 {
			return errors.New("invalid password too long")
		}
		ok, err := regexp.MatchString("^(?=.*?[a-z])(?=.*?[A-Z])(?=.*?\\d)[a-zA-Z\\d]*$", password)
		if err != nil || ok != true {
			return errors.New("invalid password pattern")
		}
	}

	if m.CreatedAt.Valid != true {
		return ErrInvalid
	} else {
		createdat := fmt.Sprint(m.CreatedAt.Time)
		ok, err := regexp.MatchString("^\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}Z$", createdat)
		if err != nil || ok != true {
			return errors.New("invalid created_at pattern")
		}
	}

	if m.UpdatedAt.Valid != true {
		return ErrInvalid
	} else {
		updatedat := fmt.Sprint(m.UpdatedAt.Time)
		ok, err := regexp.MatchString("^\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}Z$", updatedat)
		if err != nil || ok != true {
			return errors.New("invalid updated_at pattern")
		}
	}

	return nil
}

func (m PostList) Valid() error {
	return nil
}
