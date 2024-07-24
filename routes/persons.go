package routes

import (
	"eoncohub.com/person_module/models"
	"github.com/labstack/echo/v4"
)

func createPerson(context echo.Context) error {
	var person models.Person
	err := context.Bind(&person)
	if err != nil {
		return context.JSON(400, map[string]string{"message": "Invalid request"})
	}

	err = (&person).Create()
	if err != nil {
		return context.JSON(500, map[string]string{"message": err.Error()})
	}
	return context.JSON(200, map[string]string{"message": "Person created"})
}
