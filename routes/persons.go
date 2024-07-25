package routes

import (
	"eoncohub.com/person_module/models"
	"github.com/labstack/echo/v4"
	"strconv"
)

func createPerson(context echo.Context) error {
	var person models.Person
	err := context.Bind(&person)
	if err != nil {
		return context.JSON(400, map[string]string{"error": "Invalid request"})
	}

	err = (&person).Create()
	if err != nil {
		return context.JSON(500, map[string]string{"error": err.Error()})
	}
	return context.JSON(200, map[string]string{"message": "Person created"})
}

func getPerson(context echo.Context) error {
	id := context.Param("id")
	intId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return context.JSON(400, map[string]string{"error": "Invalid id"})
	}
	person, err := models.GetPerson(intId)
	if err != nil {
		return context.JSON(500, map[string]string{"error": err.Error()})
	}
	return context.JSON(200, person)
}
