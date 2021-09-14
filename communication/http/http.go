package httpserver

import (
	"strings"

	. "github.com/jinuopti/janus/log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	echoSwagger "github.com/swaggo/echo-swagger"
	"github.com/jinuopti/janus/docs"

	_ "github.com/jinuopti/janus/docs"
	"github.com/jinuopti/janus/configure"
)

// HttpServer
// @title PosGo REST API
// @version 0.1.0
// @BasePath /api/v1
// @query.collection.format multi
//
// @description <h2><b>PosGo REST API Swagger Documentation</b></h2>
//
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
//
// @tag.name Sample
// @tag.description Sample TAG
func HttpServer(port string) {
	conf := configure.GetConfig()

	port = strings.Trim(port, " ")
	if len(port) <= 0 {
		port = conf.Kredigo.HttpServerPort
	}
	docs.SwaggerInfo.Host = conf.Kredigo.SwaggerAddr + ":" + port

	e := echo.New()

	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	e.Logger.SetOutput(GetLogWriter())

	// swagger documents
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	setRoute(e)

	err := e.Start(":" + port)
	if err != nil {
		e.Logger.Fatal(err)
	}
}

// setRoute
func setRoute(e *echo.Echo) {
	// Insert API Route
	// ApiUser.SetRoute(e)
}