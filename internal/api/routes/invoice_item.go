package routes

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/dv1x3r/chinook-go/internal/api/services"

	"github.com/dv1x3r/w2go/w2"
	"github.com/labstack/echo/v5"
)

func (h *Handler) GetInvoiceItemGrid(c *echo.Context) error {
	invoiceID, err := strconv.Atoi(c.Param("invoiceID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, w2.NewErrorResponse(err.Error()))
	}
	req, err := w2.ParseGetGridRequest(c.QueryParam("request"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, w2.NewErrorResponse(err.Error()))
	}
	res, err := h.svc.GetInvoiceItemGrid(c.Request().Context(), req, invoiceID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, w2.NewErrorResponse(err.Error()))
	}
	return c.JSON(http.StatusOK, res)
}

func (h *Handler) PostInvoiceItemForm(c *echo.Context) error {
	invoiceID, err := strconv.Atoi(c.Param("invoiceID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, w2.NewErrorResponse(err.Error()))
	}
	req, err := w2.ParseSaveFormRequest[services.InvoiceItem](c.Request().Body)
	if err != nil {
		return c.JSON(http.StatusBadRequest, w2.NewErrorResponse(err.Error()))
	}
	if req.Record.ID == 0 {
		req.Record.ID, err = h.svc.CreateInvoiceItem(c.Request().Context(), req, invoiceID)
	} else {
		err = h.svc.UpdateInvoiceItem(c.Request().Context(), req, invoiceID)
	}
	if errors.Is(err, services.ErrInvoiceItemNotFound) {
		return c.JSON(http.StatusNotFound, w2.NewErrorResponse(err.Error()))
	} else if err != nil {
		return c.JSON(http.StatusInternalServerError, w2.NewErrorResponse(err.Error()))
	}
	return c.JSON(http.StatusOK, w2.NewSuccessResponse())
}

func (h *Handler) PostInvoiceItemDelete(c *echo.Context) error {
	invoiceID, err := strconv.Atoi(c.Param("invoiceID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, w2.NewErrorResponse(err.Error()))
	}
	req, err := w2.ParseRemoveGridRequest(c.Request().Body)
	if err != nil {
		return c.JSON(http.StatusBadRequest, w2.NewErrorResponse(err.Error()))
	}
	err = h.svc.DeleteInvoiceItems(c.Request().Context(), req, invoiceID)
	if errors.Is(err, services.ErrInvoiceItemNotFound) {
		return c.JSON(http.StatusNotFound, w2.NewErrorResponse(err.Error()))
	} else if err != nil {
		return c.JSON(http.StatusInternalServerError, w2.NewErrorResponse(err.Error()))
	}
	return c.JSON(http.StatusOK, w2.NewSuccessResponse())
}
