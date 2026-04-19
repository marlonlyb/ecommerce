package handlers

import (
	"errors"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/mlbautomation/Ecommmerce_MLB/domain/ports/invoice"
	"github.com/mlbautomation/Ecommmerce_MLB/infrastructure/handlers/response"
	"github.com/mlbautomation/Ecommmerce_MLB/model"
)

type Invoice struct {
	service   invoice.Service
	responser response.API
}

func NewInvoice(service invoice.Service) *Invoice {
	return &Invoice{service: service}
}

// GetByUserID returns the shops from a logged user
func (h *Invoice) GetByUserID(c echo.Context) error {
	userID, ok := c.Get("userID").(uuid.UUID)
	if !ok {
		err := model.NewError()
		err.Err = errors.New("couldn't parse ID user from token")
		err.Who = "handlers-Invoice-GetByUserID-c.Get('userID').(uuid.UUID)"
		err.APIMessage = "couldn't parse ID user from token"
		return &err
	}

	data, err := h.service.GetByUserID(userID)
	if err != nil {
		errResp := model.NewError()
		errResp.Err = err
		errResp.Who = "handlers-Invoice-MyShops-h.service.GetByUserID(userID)"
		errResp.APIMessage = err.Error()

		return &errResp
	}

	return c.JSON(h.responser.OK(data))
}

func (h *Invoice) GetAll(c echo.Context) error {
	data, err := h.service.GetAll()
	if err != nil {
		errResp := model.NewError()
		errResp.Err = err
		errResp.Who = "handlers-Invoice-GetAll-h.service.GetAll()"
		errResp.APIMessage = err.Error()

		return &errResp
	}

	return c.JSON(h.responser.OK(data))
}
