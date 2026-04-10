package handlers

import "github.com/labstack/echo/v4"

type UserHandler interface {
	Create(c echo.Context) error
	GetAll(c echo.Context) error
}

type ProductHandler interface {
	Create(c echo.Context) error
	Update(c echo.Context) error
	Delete(c echo.Context) error
	GetByID(c echo.Context) error
	GetAll(c echo.Context) error
}

type PurchaseOrderHandler interface {
	Create(c echo.Context) error
}

type LoginHandler interface {
	Login(c echo.Context) error
}

type PaypalHandler interface {
	Webhook(c echo.Context) error
}

type InvoiceHandler interface {
	GetByUserID(c echo.Context) error
	GetAll(c echo.Context) error
}
