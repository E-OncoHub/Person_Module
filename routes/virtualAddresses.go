package routes

import (
	"eoncohub.com/person_module/models"
	"github.com/labstack/echo/v4"
	"net/http"
)

func createVirtualAddress(ctx echo.Context) error {
	var virtualAddress models.VirtualAddress
	err := ctx.Bind(&virtualAddress)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]any{"error": "Invalid VirtualAddress"})
	}
	err = virtualAddress.CreateVirtualAddress()
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]any{"error": "Nu a putut crea VirtualAddress"})
	}
	return ctx.JSON(http.StatusOK, virtualAddress)
}
