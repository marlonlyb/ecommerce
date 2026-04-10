package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/mlbautomation/ProyectoEMLB/infrastructure/handlers"
)

// publicRoutes handle the routes that not requires a validation of any kind to be use
func PaypalPublic(e *echo.Echo, h handlers.PaypalHandler) {
	//recomendación de no poner "paypal" muy obvio para ataques
	route := e.Group("/api/v1/public/paypal")

	route.POST("", h.Webhook)
}
