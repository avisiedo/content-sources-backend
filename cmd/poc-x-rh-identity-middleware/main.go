package main

// https://echo.labstack.com/guide/

import (
	"encoding/json"
	"net/http"

	// m "github.com/content-services/content-sources-backend/pkg/middleware"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func logXRHI(c echo.Context) {
	identitySerialized := c.Get("x-rh-identity")
	identityJson, err := json.Marshal(identitySerialized)
	if err != nil {
		c.Logger().Infof("error x-rh-identity: %w", err)
	} else {
		c.Logger().Infof("x-rh-identity: %s", identityJson)
	}
}

func ConfigureService(e *echo.Echo) {
	if e == nil {
		return
	}

	// Setup middlewares
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))
	// e.Use(m.XRHIdentityMiddleware(m.NewXRHIdentityConfig(
	// 	func(c echo.Context) bool {
	// 		if strings.HasPrefix(c.Request().URL.Path, "/ping") {
	// 			return true
	// 		}
	// 		return false
	// 	},
	// 	m.DefaultXRHIdentityConfig.Validation,
	// )))

	// Setup routes
	e.GET("/ping", func(c echo.Context) error {
		logXRHI(c)
		return c.String(http.StatusOK, "ping")
	})
	e.GET("/hello", func(c echo.Context) error {
		logXRHI(c)
		return c.String(http.StatusOK, "Hello world")
	})

}

func main() {
	e := echo.New()
	ConfigureService(e)
	if err := e.Start("localhost:8000"); err != nil {
		panic(err)
	}
}
