package main

import (
	"embed"
	"github.com/labstack/echo/v4"
	"github.com/lyh-demo/go-webapp-demo/config"
	"github.com/lyh-demo/go-webapp-demo/container"
	"github.com/lyh-demo/go-webapp-demo/logger"
	"github.com/lyh-demo/go-webapp-demo/middleware"
	"github.com/lyh-demo/go-webapp-demo/migration"
	"github.com/lyh-demo/go-webapp-demo/repository"
	"github.com/lyh-demo/go-webapp-demo/router"
	"github.com/lyh-demo/go-webapp-demo/session"
)

//go:embed resources/config/application.*.yml
var yamlFile embed.FS

//go:embed resources/config/zaplogger.*.yml
var zapYamlFile embed.FS

//go:embed resources/public/*
var staticFile embed.FS

//go:embed resources/config/messages.properties
var propsFile embed.FS

// @title go-webapp-demo API
// @version 0.0.1
// @description This is API specification for go-webapp-demo project.

// @license.name GNU General Public License v3.0
// @license.url https://github.com/lyh-demo/go-webapp-demo/blob/main/LICENSE

// @host localhost:8080
// @BasePath /api
func main() {
	e := echo.New()

	conf, env := config.LoadAppConfig(yamlFile)
	l := logger.InitLogger(env, zapYamlFile)
	l.GetZapLogger().Infof("Loaded this configuration : application." + env + ".yml")

	messages := config.LoadMessagesConfig(propsFile)
	l.GetZapLogger().Infof("Loaded messages.properties")

	rep := repository.NewBookRepository(l, conf)
	sess := session.NewSession(l, conf)
	c := container.NewContainer(rep, sess, conf, messages, l, env)

	migration.CreateDatabase(c)
	migration.InitMasterData(c)

	router.Init(e, c)
	middleware.InitLoggerMiddleware(e, c)
	middleware.InitSessionMiddleware(e, c)
	middleware.StaticContentsMiddleware(e, c, staticFile)

	if err := e.Start(":8080"); err != nil {
		l.GetZapLogger().Errorf(err.Error())
	}

	defer rep.Close()
}
