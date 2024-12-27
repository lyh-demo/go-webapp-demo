package controller

import (
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/lyh-demo/go-webapp-demo/container"
	"net/http"
)

// APIError has an error code and a message.
type APIError struct {
	Code    int
	Message string
}

// ErrorController is a controller for handling errors.
type ErrorController interface {
	JSONError(err error, c echo.Context)
}

type errorController struct {
	container container.Container
}

// NewErrorController is constructor.
func NewErrorController(container container.Container) ErrorController {
	return &errorController{container: container}
}

// JSONError is customize error handler
func (controller *errorController) JSONError(err error, c echo.Context) {
	logger := controller.container.GetLogger()
	code := http.StatusInternalServerError
	msg := http.StatusText(code)

	var he *echo.HTTPError
	if errors.As(err, &he) {
		code = he.Code
		msg = he.Message.(string)
	}

	var apiErr APIError
	apiErr.Code = code
	apiErr.Message = msg

	if !c.Response().Committed {
		if resErr := c.JSON(code, apiErr); resErr != nil {
			logger.GetZapLogger().Errorf(resErr.Error())
		}
	}
	logger.GetZapLogger().Debugf(err.Error())
}
