package main

import (
	"eoncohub.com/person_module/db"
	"eoncohub.com/person_module/routes"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	db.InitDB()
	defer db.CloseDB()
	e := echo.New()
	e.Use(middleware.Logger())  // Log each request
	e.Use(middleware.Recover()) // Recover from panics anywhere in the chain
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"https://labstack.com", "https://labstack.net"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
	}))
	routes.RegisterRoutes(e)
	e.Logger.Fatal(e.Start(":8080"))
}
