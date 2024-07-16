package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())  // Log each request
	e.Use(middleware.Recover()) // Recover from panics anywhere in the chain

	e.Logger.Fatal(e.Start(":8080"))
}
