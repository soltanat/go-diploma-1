package entities

import "fmt"

type User struct {
	Login    Login
	Password string
	Balance  Currency
}

func NewUser(login Login, password string) *User {
	return &User{
		Login:    login,
		Password: password,
		Balance: Currency{
			Whole:   0,
			Decimal: 0,
		},
	}
}

func (u *User) Validate() error {
	if err := u.Login.Validate(); err != nil {
		return err
	}
	if len(u.Password) == 0 {
		return ValidationError{Err: fmt.Errorf("invalid password: %s", u.Password)}
	}
	return nil
}

type Login string

func (l Login) Validate() error {
	if len(l) == 0 {
		return ValidationError{Err: fmt.Errorf("invalid login: %s", l)}
	}
	return nil
}
