package handlers

import (
	"errors"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/mlbautomation/Ecommmerce_MLB/domain/ports/purchaseorder"
	"github.com/mlbautomation/Ecommmerce_MLB/infrastructure/handlers/response"
	"github.com/mlbautomation/Ecommmerce_MLB/model"
)

type PurchaseOrder struct {
	service  purchaseorder.Service
	response response.API
}

func NewPurchaseOrder(pos purchaseorder.Service) *PurchaseOrder {
	return &PurchaseOrder{service: pos}
}

func (h *PurchaseOrder) Create(c echo.Context) error {

	m := model.PurchaseOrder{}

	err := c.Bind(&m)
	if err != nil {
		return h.response.BindFailed(c, "handlers-PurchaseOrder-Create-c.Bind(&m)", err)
	}

	//Se necesita que elusuario tenga un token valido y de estos claims se saca el userID
	userID, ok := c.Get("userID").(uuid.UUID)
	if !ok {
		err = errors.New("can´t parse uuid")
		return h.response.Error(c, "handlers-PurchaseOrder-Create-c.Get('userID').(uuid.UUID)", err)
	}

	m.UserID = userID

	err = h.service.Create(&m)
	if err != nil {
		return h.response.Error(c, "handlers-PurchaseOrder-Create-h.Service.Create(&m)", err)
	}

	return c.JSON(h.response.Created(m))

}
