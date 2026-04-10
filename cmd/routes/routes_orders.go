package routes

import (
	"github.com/labstack/echo/v4"

	"github.com/mlbautomation/ProyectoEMLB/infrastructure/handlers"
)

// OrdersPrivate expone el naming oficial del MVP.
// Mantiene convivencia temporal con el flujo legado de purchaseorders.
func OrdersPrivate(e *echo.Echo, h handlers.OrderHandler, middlewares ...echo.MiddlewareFunc) {
	g := e.Group("/api/v1/private/orders", middlewares...)

	g.GET("", h.GetMine)
	g.GET("/:id", h.GetByID)
	g.POST("/checkout/paypal", h.CheckoutPayPal)
	g.POST("/:id/paypal/capture", h.CapturePayPal)
}
