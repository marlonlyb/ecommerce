package handlers

import (
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	orderport "github.com/mlbautomation/ProyectoEMLB/domain/ports/order"
	"github.com/mlbautomation/ProyectoEMLB/domain/services"
	"github.com/mlbautomation/ProyectoEMLB/infrastructure/handlers/response"
	"github.com/mlbautomation/ProyectoEMLB/model"
)

type Order struct {
	service orderport.Service
}

func NewOrder(service orderport.Service) *Order {
	return &Order{service: service}
}

func (h *Order) CheckoutPayPal(c echo.Context) error {
	userID, err := parseUserID(c)
	if err != nil {
		return response.ContractError(401, "authentication_required", "Debes iniciar sesión para continuar")
	}

	request := model.CheckoutPayPalRequest{}
	if err = c.Bind(&request); err != nil {
		return response.ContractError(400, "validation_error", "Los datos enviados no son válidos")
	}

	orderData, payPalOrderID, err := h.service.CheckoutPayPal(userID, request)
	if err != nil {
		return mapOrderError(err)
	}

	return c.JSON(response.ContractCreated(map[string]interface{}{
		"order":  orderData,
		"paypal": map[string]string{"order_id": payPalOrderID},
	}))
}

func (h *Order) CapturePayPal(c echo.Context) error {
	userID, err := parseUserID(c)
	if err != nil {
		return response.ContractError(401, "authentication_required", "Debes iniciar sesión para continuar")
	}

	orderID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.ContractError(400, "validation_error", "El identificador de la orden no es válido")
	}

	request := model.CapturePayPalOrderRequest{}
	if err = c.Bind(&request); err != nil {
		return response.ContractError(400, "validation_error", "Los datos enviados no son válidos")
	}

	orderData, err := h.service.CapturePayPal(userID, orderID, request)
	if err != nil {
		return mapOrderError(err)
	}

	return c.JSON(response.ContractOK(map[string]interface{}{"order": orderData}))
}

func (h *Order) GetMine(c echo.Context) error {
	userID, err := parseUserID(c)
	if err != nil {
		return response.ContractError(401, "authentication_required", "Debes iniciar sesión para continuar")
	}

	orders, err := h.service.ListByUserID(userID)
	if err != nil {
		return response.ContractError(500, "unexpected_error", "No fue posible obtener las órdenes")
	}

	return c.JSON(response.ContractOK(map[string]interface{}{"items": orders}))
}

func (h *Order) GetByID(c echo.Context) error {
	userID, err := parseUserID(c)
	if err != nil {
		return response.ContractError(401, "authentication_required", "Debes iniciar sesión para continuar")
	}

	orderID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.ContractError(400, "validation_error", "El identificador de la orden no es válido")
	}

	orderData, err := h.service.GetByIDForUser(orderID, userID)
	if err != nil {
		return mapOrderError(err)
	}

	return c.JSON(response.ContractOK(orderData))
}

func parseUserID(c echo.Context) (uuid.UUID, error) {
	userID, ok := c.Get("userID").(uuid.UUID)
	if !ok || userID == uuid.Nil {
		return uuid.Nil, errors.New("invalid user")
	}

	return userID, nil
}

func mapOrderError(err error) error {
	if errors.Is(err, services.ErrValidation) {
		if strings.Contains(strings.ToLower(err.Error()), "inactive") {
			return response.ContractError(409, "product_inactive", "Uno o más productos ya no están disponibles")
		}
		return response.ContractError(400, "validation_error", "Los datos enviados no son válidos")
	}

	if errors.Is(err, services.ErrStockInsufficient) {
		return response.ContractError(409, "stock_insufficient", "Una o más variantes ya no tienen stock suficiente")
	}

	if errors.Is(err, services.ErrOrderNotFound) {
		return response.ContractError(404, "not_found", "Orden no encontrada")
	}

	if errors.Is(err, services.ErrOrderStateInvalid) {
		return response.ContractError(409, "order_state_invalid", "La orden no puede capturarse en el estado actual")
	}

	if errors.Is(err, services.ErrPayPalCaptureFailed) {
		return response.ContractError(422, "paypal_capture_failed", "PayPal no confirmó la captura")
	}

	if strings.Contains(strings.ToLower(err.Error()), "insufficient stock") {
		return response.ContractError(409, "stock_insufficient", "Una o más variantes ya no tienen stock suficiente")
	}

	return response.ContractError(500, "unexpected_error", "No fue posible procesar la orden")
}
