package test

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/lyh-demo/go-webapp-demo/config"
	"github.com/lyh-demo/go-webapp-demo/container"
	"github.com/lyh-demo/go-webapp-demo/logger"
	"github.com/lyh-demo/go-webapp-demo/middleware"
	"github.com/lyh-demo/go-webapp-demo/migration"
	"github.com/lyh-demo/go-webapp-demo/repository"
	"github.com/lyh-demo/go-webapp-demo/session"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
	"net/http"
	"net/http/httptest"
	"strings"
)

// PrepareForControllerTest func prepares the controllers for testing.
func PrepareForControllerTest(isSecurity bool) (*echo.Echo, container.Container) {
	e := echo.New()

	conf := createConfig(isSecurity)
	l := initTestLogger()
	c := initContainer(conf, l)

	middleware.InitLoggerMiddleware(e, c)

	migration.CreateDatabase(c)
	migration.InitMasterData(c)

	middleware.InitSessionMiddleware(e, c)
	return e, c
}

// PrepareForServiceTest func prepares the services for testing.
func PrepareForServiceTest() container.Container {
	conf := createConfig(false)
	l := initTestLogger()
	c := initContainer(conf, l)

	migration.CreateDatabase(c)
	migration.InitMasterData(c)

	return c
}

// PrepareForLoggerTest func prepares the loggers for testing.
func PrepareForLoggerTest() (*echo.Echo, container.Container, *observer.ObservedLogs) {
	e := echo.New()

	conf := createConfig(false)
	l, observedLogs := initObservedLogger()
	c := initContainer(conf, l)

	migration.CreateDatabase(c)
	migration.InitMasterData(c)

	middleware.InitSessionMiddleware(e, c)
	middleware.InitLoggerMiddleware(e, c)
	return e, c, observedLogs
}

func createConfig(isSecurity bool) *config.Config {
	conf := &config.Config{}
	conf.Database.Dialect = "sqlite3"
	conf.Database.Host = "file::memory:?cache=shared"
	conf.Database.Migration = true
	conf.Extension.MasterGenerator = true
	conf.Extension.SecurityEnabled = isSecurity
	conf.Log.RequestLogFormat = "${remote_ip} ${account_name} ${uri} ${method} ${status}"
	return conf
}

func initContainer(conf *config.Config, logger logger.Logger) container.Container {
	rep := repository.NewBookRepository(logger, conf)
	sess := session.NewSession(logger, conf)
	messages := map[string]string{
		"ValidationErrMessageBookTitle": "Please enter the title with 3 to 50 characters.",
		"ValidationErrMessageBookISBN":  "Please enter the ISBN with 10 to 20 characters."}
	c := container.NewContainer(rep, sess, conf, messages, logger, "test")
	return c
}

func initTestLogger() logger.Logger {
	myConfig := createLoggerConfig()
	z, err := myConfig.Build()
	if err != nil {
		fmt.Printf("Error")
	}
	sugar := z.Sugar()
	// set package variable logger.
	l := logger.NewLogger(sugar)
	l.GetZapLogger().Infof("Success to read zap logger configuration")
	_ = z.Sync()
	return l
}

func initObservedLogger() (logger.Logger, *observer.ObservedLogs) {
	observedZapCore, observedLogs := observer.New(zap.DebugLevel)
	sugar := zap.New(observedZapCore).Sugar()

	// set package variable logger.
	l := logger.NewLogger(sugar)
	return l, observedLogs
}

func createLoggerConfig() zap.Config {
	level := zap.NewAtomicLevel()
	level.SetLevel(zapcore.DebugLevel)

	return zap.Config{
		Level:       level,
		Encoding:    "console",
		Development: true,
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "Time",
			LevelKey:       "Level",
			NameKey:        "Name",
			CallerKey:      "Caller",
			MessageKey:     "Msg",
			StacktraceKey:  "St",
			EncodeLevel:    zapcore.CapitalLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}
}

// ConvertToString func converts model to string.
func ConvertToString(model interface{}) string {
	bytes, _ := json.Marshal(model)
	return string(bytes)
}

// NewJSONRequest func creates a new request using JSON format.
func NewJSONRequest(method string, target string, param interface{}) *http.Request {
	req := httptest.NewRequest(method, target, strings.NewReader(ConvertToString(param)))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	return req
}

// GetCookie func gets a cookie from an HTTP request.
func GetCookie(rec *httptest.ResponseRecorder, cookieName string) string {
	parser := &http.Request{Header: http.Header{"Cookie": rec.Header()["Set-Cookie"]}}
	if cookie, err := parser.Cookie(cookieName); cookie != nil && err == nil {
		return cookie.Value
	}
	return ""
}
