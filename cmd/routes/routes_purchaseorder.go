package routes

import (
	"github.com/labstack/echo/v4"

	"github.com/mlbautomation/Ecommmerce_MLB/infrastructure/handlers"
)

func PurchaseOrderPrivate(e *echo.Echo, h handlers.PurchaseOrderHandler, middlewares ...echo.MiddlewareFunc) {
	g := e.Group("/api/v1/private/purchaseorders", middlewares...)

	g.POST("", h.Create)
}
