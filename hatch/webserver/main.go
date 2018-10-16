package main

import (
	"aria/hatch/webserver/handler"
	"aria/hatch/webserver/middleware"
	"aria/hatch/webserver/service"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/swag/gen"
)

// @title Aria WebServer API
// @version 1.0
// @host localhost:8080
func main() {
	// init config
	if err := InitConfig("./config/config.yaml"); err != nil {
		panic(err)
	}

	// init microservice
	if err := service.InitService(Config().Services); err != nil {
		panic(err)
	}

	// start web server
	g := gin.New()
	g.Use(gin.Recovery())
	g.Use(middleware.Logger())
	g.Use(cors.Default())
	g.Use(middleware.ZipkinTracing(
		Config().Statistic.Tracing.Zipkin.Url,
		"webserver",
		fmt.Sprintf("0.0.0.0:%d", Config().Server.Port),
		Config().Statistic.Enable,
	))

	authorizedGroup := g.Group("/")
	authorizedGroup.GET("/production/get/:id", handler.GetProduction)
	authorizedGroup.POST("/production/add", handler.AddProduction)

	// process swagger api docs
	swagger(g)

	address := fmt.Sprintf(":%d", Config().Port)
	if err := g.Run(address); err != nil {
		panic(err)
	}
}

func swagger(engine *gin.Engine) {
	if Config().GenerateDocs {
		gen.New().Build("./", "main.go", "./swagger", "camelcase")
	}
	if Config().SwaggerServerOn {
		engine.Static("/swagger", "./swagger")
	}
}
