package routes

import "github.com/labstack/echo/v4"

func RegisterRoutes(server *echo.Echo) {
	server.POST("/virtual_address", createVirtualAddress)
	server.POST("/person/create", createPerson)
}
