package handlers

import (
	"github.com/labstack/echo/v4"

	"github.com/mlbautomation/ProyectoEMLB/domain/ports/user"
	"github.com/mlbautomation/ProyectoEMLB/infrastructure/handlers/response"
	"github.com/mlbautomation/ProyectoEMLB/model"
)

/* date cuenta que a diferencia de model
aqui los tipos y funciones son privadas
porque serán consultadas desde la ruta */

type User struct {
	service   user.Service
	responser response.API
}

func NewUser(us user.Service) *User {
	return &User{service: us}
}

func (h *User) Create(c echo.Context) error {
	m := model.User{}

	//vinculamos (bind) la información del cuerpo de la solicitud
	err := c.Bind(&m)
	if err != nil {
		return h.responser.BindFailed(c, "handlers-User-Create-c.Bind(&m)", err)
		//return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	err = h.service.Create(&m)
	if err != nil {
		return h.responser.Error(c, "handlers-User-Create-h.service.Create((&m))", err)
		//return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(h.responser.Created(m))
	//return c.JSON(http.StatusOK, m)
}

func (h *User) GetAll(c echo.Context) error {
	users, err := h.service.GetAll()
	if err != nil {
		return h.responser.Error(c, "handlers-User-Create-h.service.GetAll()", err)
	}

	return c.JSON(h.responser.OK(users))
}
