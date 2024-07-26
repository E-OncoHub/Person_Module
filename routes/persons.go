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

// In persons.go
func updatePerson(context echo.Context) error {
	// Parse the person ID from the URL
	id := context.Param("id")
	personID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return context.JSON(400, map[string]string{"error": "Invalid person ID"})
	}

	// Bind the updated person data from the request body
	var updatedPerson models.Person
	if err := context.Bind(&updatedPerson); err != nil {
		return context.JSON(400, map[string]string{"error": "Invalid request body"})
	}

	// Set the ID of the person to update
	updatedPerson.IDPerson = personID

	// Call the Update method on the person model
	err = updatedPerson.Update()
	if err != nil {
		return context.JSON(500, map[string]string{"error": err.Error()})
	}

	return context.JSON(200, map[string]string{"message": "Person updated successfully"})
}
