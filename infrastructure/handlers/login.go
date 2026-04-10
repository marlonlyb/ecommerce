package handlers

import (
	"os"
	"strings"

	"github.com/labstack/echo/v4"

	"github.com/mlbautomation/ProyectoEMLB/domain/ports/login"
	"github.com/mlbautomation/ProyectoEMLB/infrastructure/handlers/response"
)

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Login struct {
	service   login.Service
	responser response.API
}

func NewLogin(us login.Service) *Login {
	return &Login{service: us}
}

func (h *Login) Login(c echo.Context) error {

	m := loginRequest{}

	err := c.Bind(&m)
	if err != nil {
		return h.responser.BindFailed(c, "handlers-Login-Login-c.Bind(&m)", err)
	}

	userModel, tokenSigned, err := h.service.Login(m.Email, m.Password, os.Getenv("JWT_SECRET_KEY"))
	if err != nil {
		if strings.Contains(err.Error(), "crypto/bcrypt: hashedPassword is not the hash of the given password") ||
			strings.Contains(err.Error(), "no rows in result set") {
			return h.responser.HashedPassword(c, "handlers-Login-Login-h.service.Login()", err)
		}
		return h.responser.Error(c, "handlers-Login-Login-h.service.Login()", err)
	}

	data := map[string]interface{}{"user": userModel, "token": tokenSigned}
	return c.JSON(h.responser.OK(data))
}
