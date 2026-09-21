package routes

import (
	"errors"
	"net/http"

	"github.com/dv1x3r/chinook-go/internal/api/services"

	"github.com/dv1x3r/w2go/w2"
	"github.com/labstack/echo/v5"
)

func (h *Handler) GetEmployeeGrid(c *echo.Context) error {
	req, err := w2.ParseGetGridRequest(c.QueryParam("request"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, w2.NewErrorResponse(err.Error()))
	}
	res, err := h.svc.GetEmployeeGrid(c.Request().Context(), req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, w2.NewErrorResponse(err.Error()))
	}
	return c.JSON(http.StatusOK, res)
}

func (h *Handler) GetEmployeeDropdown(c *echo.Context) error {
	req, err := w2.ParseGetDropdownRequest(c.QueryParam("request"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, w2.NewErrorResponse(err.Error()))
	}
	res, err := h.svc.GetEmployeeDropdown(c.Request().Context(), req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, w2.NewErrorResponse(err.Error()))
	}
	return c.JSON(http.StatusOK, res)
}

func (h *Handler) PostEmployeeForm(c *echo.Context) error {
	req, err := w2.ParseSaveFormRequest[services.Employee](c.Request().Body)
	if err != nil {
		return c.JSON(http.StatusBadRequest, w2.NewErrorResponse(err.Error()))
	}
	if req.Record.ID == 0 {
		req.Record.ID, err = h.svc.CreateEmployee(c.Request().Context(), req)
	} else {
		err = h.svc.UpdateEmployee(c.Request().Context(), req)
	}
	if errors.Is(err, services.ErrEmployeeNotFound) {
		return c.JSON(http.StatusNotFound, w2.NewErrorResponse(err.Error()))
	} else if err != nil {
		return c.JSON(http.StatusInternalServerError, w2.NewErrorResponse(err.Error()))
	}
	return c.JSON(http.StatusOK, w2.NewSuccessResponse())
}

func (h *Handler) PostEmployeeDelete(c *echo.Context) error {
	req, err := w2.ParseRemoveGridRequest(c.Request().Body)
	if err != nil {
		return c.JSON(http.StatusBadRequest, w2.NewErrorResponse(err.Error()))
	}
	err = h.svc.DeleteEmployees(c.Request().Context(), req)
	if errors.Is(err, services.ErrEmployeeHasRelated) {
		return c.JSON(http.StatusConflict, w2.NewErrorResponse(err.Error()))
	} else if errors.Is(err, services.ErrEmployeeNotFound) {
		return c.JSON(http.StatusNotFound, w2.NewErrorResponse(err.Error()))
	} else if err != nil {
		return c.JSON(http.StatusInternalServerError, w2.NewErrorResponse(err.Error()))
	}
	return c.JSON(http.StatusOK, w2.NewSuccessResponse())
}
