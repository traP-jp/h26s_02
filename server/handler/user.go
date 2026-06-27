package handler

import "github.com/labstack/echo/v5"

type User struct {
}

func NewUser() *User {
	return &User{}
}

func (u *User) GetMe(_ *echo.Context) error {
	// TODO
	return nil
}
