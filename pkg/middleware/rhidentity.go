package middleware

import (
	"encoding/base64"
	"encoding/json"
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/redhatinsights/platform-go-middlewares/identity"
)

const XRHIdentityKey = "x-rh-identity"

type XRHIdentityValidation func(value *identity.XRHID) bool

type XRHIdentityConfig struct {
	// Skipper defines a function to skip middleware.
	Skipper middleware.Skipper
	// Validation defines a function to validate the XRHID content
	Validation XRHIdentityValidation
}

var DefaultXRHIdentityConfig XRHIdentityConfig = XRHIdentityConfig{
	Skipper:    middleware.DefaultSkipper,
	Validation: DefaultXRHIdentityValidation,
}

func DefaultXRHIdentityValidation(value *identity.XRHID) bool {
	if value == nil {
		return false
	}
	return true
}

func NewXRHIdentityConfig(s middleware.Skipper, v XRHIdentityValidation) XRHIdentityConfig {
	if s == nil {
		s = DefaultXRHIdentityConfig.Skipper
	}
	if v == nil {
		v = DefaultXRHIdentityConfig.Validation
	}
	return XRHIdentityConfig{
		Skipper:    s,
		Validation: v,
	}
}

func XRHIdentityMiddleware(config XRHIdentityConfig) echo.MiddlewareFunc {
	if config.Skipper == nil {
		config.Skipper = DefaultXRHIdentityConfig.Skipper
	}
	if config.Validation == nil {
		config.Validation = DefaultXRHIdentityConfig.Validation
	}
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Check skipper
			if config.Skipper(c) {
				return next(c)
			}

			var (
				decodedIdentity   []byte
				err               error
				identityHeader    identity.XRHID
				nodecodedIdentity string
			)

			// Decode header
			if nodecodedIdentity = c.Request().Header.Get(XRHIdentityKey); nodecodedIdentity == "" {
				c.Logger().Errorf("Error '%s' header is empty", XRHIdentityKey, nodecodedIdentity)
				return echo.NewHTTPError(http.StatusUnauthorized, "")
			}
			if decodedIdentity, err = base64.StdEncoding.DecodeString(nodecodedIdentity); err != nil {
				c.Logger().Errorf("Error decoding '%s' header: %w", XRHIdentityKey, err)
				return echo.NewHTTPError(http.StatusUnauthorized, "")
			}

			// Parse header
			if err = json.Unmarshal(decodedIdentity, &identityHeader); err != nil {
				c.Logger().Errorf("Error unserializing '%s' header: %w", XRHIdentityKey, err)
				return echo.NewHTTPError(http.StatusUnauthorized, "")
			}

			// Validate the current content
			if config.Validation != nil {
				if !config.Validation(&identityHeader) {
					c.Logger().Errorf("Error validating '%s' header: %w", XRHIdentityKey, err)
					return echo.NewHTTPError(http.StatusUnauthorized, "")
				}
			}

			// Set validated value into the context
			c.Set(XRHIdentityKey, identityHeader)

			// Call next middleware
			return next(c)
		}
	}
}
