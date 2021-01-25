package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"insectt.io/api/websocket"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Static("/", "../public")
	e.GET("/ws", websocket.DbHandler)
	e.Logger.Fatal(e.Start(":8080"))
}
