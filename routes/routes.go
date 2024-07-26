package routes

import "github.com/labstack/echo/v4"

func RegisterRoutes(server *echo.Echo) {
	// Virtual Address routes
	server.POST("/virtual_address", createVirtualAddress)

	// Person routes
	server.POST("/person/create", createPerson)
	server.GET("/person/:id", getPerson)
	server.PUT("/person/update/:id", updatePerson)
	server.DELETE("/person/delete/:id", deletePerson)

}
