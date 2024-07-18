package main

import (
	"github.com/cruseraca/json-to-table/handlers"
	"github.com/labstack/echo/v4"
)

func main()  {
	e := echo.New()

	checkJsonHandler := handlers.NewCheckJsonHandler()

	v1_0 := e.Group("/api/v1.0")
	v1_0.POST("/check-json", checkJsonHandler.CheckJson)
	v1_0.POST("/generate-table", checkJsonHandler.GenerateTable)
	v1_0.POST("/generate-table-order", checkJsonHandler.GenerateTableOrder)
	
	e.Logger.Fatal(e.Start(":1323"))
}