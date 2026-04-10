package handlers

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/mlbautomation/ProyectoEMLB/domain/ports/product"
	"github.com/mlbautomation/ProyectoEMLB/infrastructure/handlers/response"
	"github.com/mlbautomation/ProyectoEMLB/model"
)

type Product struct {
	service   product.Service
	responser response.API
}

func NewProduct(ps product.Service) *Product {
	return &Product{service: ps}
}

func (h *Product) Create(c echo.Context) error {
	m := model.Product{}

	err := c.Bind(&m)
	if err != nil {
		return h.responser.BindFailed(c, "handlers-Product-Create-c.Bind(&m)", err)
	}

	err = h.service.Create(&m)
	if err != nil {
		return h.responser.Error(c, "handlers-Product-Create-h.service.Create()", err)
	}

	return c.JSON(h.responser.Created(m))
}

func (h *Product) Update(c echo.Context) error {
	m := model.Product{}

	err := c.Bind(&m)
	if err != nil {
		return h.responser.BindFailed(c, "handlers-Product-Update-c.Bind(&m)", err)
	}

	ID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return h.responser.Error(c, "handlers-Product-Update-uuid.Parse(c.Param('id'))", err)
	}
	m.ID = ID

	err = h.service.Update(&m)
	if err != nil {
		return h.responser.Error(c, "handlers-Product-Update-h.service.Update(&m)", err)
	}

	return c.JSON(h.responser.Updated(m))
}

func (h *Product) Delete(c echo.Context) error {
	ID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return h.responser.Error(c, "handlers-Product-Delete-uuid.Parse(c.Param('id'))", err)
	}

	err = h.service.Delete(ID)
	if err != nil {
		return h.responser.Error(c, "handlers-Product-Delete-h.service.Delete(ID)", err)
	}

	return c.JSON(h.responser.Deleted(nil))
}

func (h *Product) GetByID(c echo.Context) error {
	ID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return h.responser.Error(c, "handlers-Product-GetByID-uuid.Parse(c.Param('id'))", err)
	}

	productData, err := h.service.GetByID(ID)
	if err != nil {
		return h.responser.Error(c, "handlers-Product-GetByID-h.service.GetByID(ID)", err)
	}

	return c.JSON(h.responser.OK(productData))
}

/* Paginar el GETALL, en el query param del endpoint recibimos:
limit(cuantos registros quieren recibir) y page (en que páquina quieren mostrar)
offset: se genera limit*pag -limit */

func (h *Product) GetAll(c echo.Context) error {
	products, err := h.service.GetAll()
	if err != nil {
		return h.responser.Error(c, "handlers-Product-GetAll-h.service.GetAll()", err)
	}

	return c.JSON(h.responser.OK(products))
}
