package handler

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

type User struct {
}

func NewUser() *User {
	return &User{}
}

type GetMeResponse struct {
	UserName string `json:"userName"`
}

func (u *User) GetMe(c *echo.Context) error {
	userName, err := GetUserName(c)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, GetMeResponse{
		UserName: userName,
	})
}
