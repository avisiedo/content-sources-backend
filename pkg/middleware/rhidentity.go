package middleware

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/labstack/echo/v4"
)

type (
	RedhatIdentityConfig struct {
		// Skipper Skipper
	}
	RedhatIdentityHeader struct {
		Identity string `header:"x-rh-identity"`
	}
	RedHatIdentity struct {
		AccountNumber  string `json:"identity.account_number"`
		OrganizationId string `json:"identity.internal.org_id"`
	}
)

var DefaultRedhatIdentityConfig = RedhatIdentityConfig{}

func (r *RedHatIdentity) Validate() error {
	if r.AccountNumber == "" {
		return fmt.Errorf("AccountNumber is empty")
	}
	if r.OrganizationId == "" {
		return fmt.Errorf("OrganizationId is empty")
	}
	return nil
}

// // https://github.com/labstack/echo/issues/514
// // https://github.com/labstack/echo/blob/master/middleware/logger.go
// // https://stackoverflow.com/questions/69326129/does-set-method-of-echo-context-saves-the-value-to-the-underlying-context-cont
// // How to set in the request context the values unmarshalled from 'x-rh-identity' header
// func RedhatIdentityMiddleware(fn echo.HandlerFunc) echo.MiddlewareFunc {
// 	return func(ctx echo.Context) error {
// 		context.Context
// 		ctx.Request().WithContext()
// 		ctx.Request().WithContext()

// 		return fn(NewRHIContext(ctx))
// 	}
// }
func RedhatIdentityMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			var rhih RedhatIdentityHeader
			var err error

			if err = c.Bind(&rhih); err != nil {
				return err
			}

			var rhihDecoded []byte
			if rhihDecoded, err = base64.StdEncoding.DecodeString(rhih.Identity); err != nil {
				return err
			}

			var rhIdentity RedHatIdentity
			if err = json.Unmarshal(rhihDecoded, &rhIdentity); err != nil {
				return err
			}

			if err = rhIdentity.Validate(); err != nil {
				return err
			}
			c.Set("x-rh-identity.account-number", rhIdentity.AccountNumber)
			c.Set("x-rh-identity.organization-id", rhIdentity.OrganizationId)
			echo.Logger.Info("RedhatIdentityMiddleware", "x-rh-identity.account-number", rhIdentity.AccountNumber)
			echo.Logger.Info("RedhatIdentityMiddleware", "x-rh-identity.organization-id", rhIdentity.OrganizationId)
			return next(c)
		}
	}
}

// func NewRHIContext(ctx echo.Context) RedhatIdentityMiddlewareContext {
// 	return RedhatIdentityMiddlewareContext{
// 		OldContext: ctx,
// 	}
// }

// type contextValue struct {
// 	echo.Context
// }

// // Get retrieves data from the context.
// func (ctx contextValue) Get(key string) interface{} {
// 	// get old context value
// 	val := ctx.Context.Get(key)
// 	if val != nil {
// 		return val
// 	}
// 	return ctx.Request().Context().Value(key)
// }

// // Get retrieves data from the context.
// func (ctx contextValue) Get(key string) interface{} {
// 	// get old context value
// 	val := ctx.Context.Get(key)
// 	if val != nil {
// 		return val
// 	}
// 	return ctx.Request().Context().Value(key)
// }

// // Set saves data in the context.
// func (ctx contextValue) Set(key string, val interface{}) {
// 	ctx.SetRequest(ctx.Request().WithContext(context.WithValue(ctx.Request().Context(), key, val)))
// }
