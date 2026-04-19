package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/mlbautomation/Ecommmerce_MLB/infrastructure/handlers"
)

func LoginPublic(e *echo.Echo, h handlers.LoginHandler) {
	g := e.Group("/api/v1/public/login")

	g.POST("", h.Login)
}
