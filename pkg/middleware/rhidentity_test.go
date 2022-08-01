package middleware

import (
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/labstack/echo"
	"github.com/redhatinsights/platform-go-middlewares/identity"
	"github.com/stretchr/testify/assert"
)

func TestDefaultXRHIdentityValidation(t *testing.T) {
	var (
		sut    XRHIdentityValidation
		id     identity.XRHID
		result bool
	)
	sut = DefaultXRHIdentityValidation

	result = sut(nil)
	assert.False(t, result)

	id = identity.XRHID{
		Identity: identity.Identity{
			OrgID: "12345",
			Type:  "Associate",
		},
	}
	result = sut(&id)
	assert.True(t, result)
}

func mySkipper(c echo.Context) bool {
	return false
}

func myValidator(identity *identity.XRHID) bool {
	if identity == nil {
		return false
	}
	return true
}

func TestNewXRHIdentityConfig(t *testing.T) {
	var (
		config XRHIdentityConfig
	)
	config = NewXRHIdentityConfig(nil, nil)
	assert.NotNil(t, config.Skipper)
	assert.Equal(t, reflect.ValueOf(DefaultXRHIdentityConfig.Skipper).Pointer(), reflect.ValueOf(config.Skipper).Pointer())
	assert.NotNil(t, config.Validation)
	assert.Equal(t, reflect.ValueOf(DefaultXRHIdentityConfig.Validation).Pointer(), reflect.ValueOf(config.Validation).Pointer())

	config = NewXRHIdentityConfig(mySkipper, nil)
	assert.NotNil(t, config.Skipper)
	assert.Equal(t, reflect.ValueOf(mySkipper).Pointer(), reflect.ValueOf(config.Skipper).Pointer())
	assert.NotNil(t, config.Validation)
	assert.Equal(t, reflect.ValueOf(DefaultXRHIdentityConfig.Validation).Pointer(), reflect.ValueOf(config.Validation).Pointer())

	config = NewXRHIdentityConfig(nil, myValidator)
	assert.NotNil(t, config.Skipper)
	assert.Equal(t, reflect.ValueOf(DefaultXRHIdentityConfig.Skipper).Pointer(), reflect.ValueOf(config.Skipper).Pointer())
	assert.NotNil(t, config.Validation)
	assert.Equal(t, reflect.ValueOf(myValidator).Pointer(), reflect.ValueOf(config.Validation).Pointer())

	config = NewXRHIdentityConfig(mySkipper, myValidator)
	assert.NotNil(t, config.Skipper)
	assert.Equal(t, reflect.ValueOf(mySkipper).Pointer(), reflect.ValueOf(config.Skipper).Pointer())
	assert.NotNil(t, config.Validation)
	assert.Equal(t, reflect.ValueOf(myValidator).Pointer(), reflect.ValueOf(config.Validation).Pointer())
}

func TestXRHIdentityMiddleware(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()
	c := e.NewContext(req, res)

	h := XRHIdentityMiddleware(DefaultXRHIdentityConfig)(func(c echo.Context) error {
		return c.String(http.StatusOK, "test")
	})

	assert := assert.New(t)
	validJson := []string{
		`{ "identity": {"account_number": "540155", "auth_type": "jwt-auth", "org_id": "1979710", "type": "User", "internal": {"org_id": "1979710"} } }`,
		`{ "identity": {"account_number": "540155", "auth_type": "cert-auth", "org_id": "1979710", "type": "Associate", "internal": {"org_id": "1979710"} } }`,
		`{ "identity": {"account_number": "540155", "auth_type": "basic-auth", "type": "Associate", "internal": {"org_id": "1979710"} } }`,
		`{ "identity": {"account_number": "540155", "auth_type": "cert-auth", "org_id": "1979710", "type": "Associate", "internal": {} } }`,
	}
	for _, json := range validJson {
		data := base64.StdEncoding.EncodeToString([]byte(json))
		req.Header.Set(XRHIdentityKey, data)
		assert.NoError(h(c))
	}
}
