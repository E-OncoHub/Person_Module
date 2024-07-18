package routes

import (
	"eoncohub.com/person_module/models"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

func getJudById(ctx echo.Context) error {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]any{"error": "Invalid ID"})
	}
	jud, err := models.GetJudById(id)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]any{"error": "Nu a gasit Judetul"})
	}
	return ctx.JSON(http.StatusOK, jud)
}
