package routes

import "github.com/labstack/echo/v4"

func RegisterRoutes(server *echo.Echo) {
	server.GET("/address/jud/:id", getJudById)
}
