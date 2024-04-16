package main

import (
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	e.Static("/", "./static/")
	e.GET("/", func(c echo.Context) error {
		return Overview().Render(c.Request().Context(), c.Response().Writer)
	})
	e.Logger.Fatal(e.Start(":1323"))
}
