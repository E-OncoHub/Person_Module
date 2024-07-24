package routes

import (
	"eoncohub.com/person_module/db"
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
	tx, err := db.DB.Begin()
	err = virtualAddress.CreateVirtualAddress(tx)
	if err != nil {
		tx.Rollback()
		return ctx.JSON(http.StatusInternalServerError, map[string]any{"error": "Nu a putut crea VirtualAddress"})
	}
	err = tx.Commit()
	return ctx.JSON(http.StatusOK, virtualAddress)
}
