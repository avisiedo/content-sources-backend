package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/content-services/content-sources-backend/pkg/client"
	"github.com/content-services/content-sources-backend/pkg/config"
	mocks_client "github.com/content-services/content-sources-backend/pkg/test/mocks/client"
	"github.com/labstack/echo/v4"
	"github.com/redhatinsights/platform-go-middlewares/identity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFromHttpVerbToRbacVerb(t *testing.T) {
	type TestCase struct {
		Name     string
		Given    string
		Expected client.RbacVerb
	}

	testCases := []TestCase{
		{
			Name:     "empty method",
			Given:    "",
			Expected: client.RbacVerbUndefined,
		},
		{
			Name:     "non existing method",
			Given:    "ANYOTHERTHING",
			Expected: client.RbacVerbUndefined,
		},
		{
			Name:     "GET",
			Given:    echo.GET,
			Expected: client.RbacVerbRead,
		},
		{
			Name:     "POST",
			Given:    echo.POST,
			Expected: client.RbacVerbWrite,
		},
		{
			Name:     "PUT",
			Given:    echo.PUT,
			Expected: client.RbacVerbWrite,
		},
		{
			Name:     "PATCH",
			Given:    echo.PATCH,
			Expected: client.RbacVerbWrite,
		},
		{
			Name:     "DELETE",
			Given:    echo.DELETE,
			Expected: client.RbacVerbWrite,
		},
	}

	for _, testCase := range testCases {
		t.Log(testCase.Name)
		result := fromHttpVerbToRbacVerb(testCase.Given)
		assert.Equal(t, testCase.Expected, result)
	}
}

func TestFromPathToResource(t *testing.T) {

}

func mockXRhUserIdentity(t *testing.T, org_id string, accNumber string) string {
	var (
		err       error
		xrhid     identity.XRHID
		jsonBytes []byte
	)
	xrhid.Identity.OrgID = org_id
	xrhid.Identity.AccountNumber = accNumber
	xrhid.Identity.Internal.OrgID = org_id

	jsonBytes, err = json.Marshal(xrhid)
	require.NoError(t, err)

	return string(jsonBytes)
}

func rbacServe(t *testing.T, req *http.Request, resource string, verb client.RbacVerb, rbacAllowed bool, rbacError error) *httptest.ResponseRecorder {
	var (
		xrhid string
		rw    *httptest.ResponseRecorder
	)

	require.NotNil(t, req)

	xrhid = mockXRhUserIdentity(t, "12345", "12345")
	require.NotEqual(t, "", xrhid)

	mockRbacClient := mocks_client.NewRbac(t)
	require.NotNil(t, mockRbacClient)
	mockRbacClient.On("Allowed", xrhid, resource, verb).Return(rbacAllowed, rbacError)

	e := echo.New()
	require.NotNil(t, e)

	e.Use(
		// Add Rbac middleware
		NewRbac(Rbac{
			BaseUrl: config.Get().Clients.RbacBaseUrl,
		}, mockRbacClient),
	)

	rw = httptest.NewRecorder()
	require.NotNil(t, rw)

	req.Header.Set(xrhidHeader, xrhid)

	e.ServeHTTP(rw, req)
	mockRbacClient.AssertExpectations(t)

	return rw
}

func TestRbacMiddleware(t *testing.T) {
	type TestCaseGivenRbac struct {
		Resource string
		Verb     client.RbacVerb
		Allowed  bool
		Err      error
	}
	type TestCaseGivenRequest struct {
		Method string
		Path   string
	}
	type TestCaseGiven struct {
		Request      TestCaseGivenRequest
		MockResponse TestCaseGivenRbac
	}
	type TestCaseExpected struct {
		Code int
		Body string
	}
	type TestCase struct {
		Name     string
		Given    TestCaseGiven
		Expected TestCaseExpected
	}

	testCases := []TestCase{
		{
			Name: "TODO Update: simple test",
			Given: TestCaseGiven{
				Request: TestCaseGivenRequest{
					Method: http.MethodGet,
					Path:   "/api/content-sources/repositories/",
				},
				MockResponse: TestCaseGivenRbac{
					Resource: "repositories",
					Verb:     client.RbacVerbRead,
					Allowed:  true,
					Err:      nil,
				},
			},
			Expected: TestCaseExpected{
				Code: http.StatusOK,
				Body: ``,
			},
		},
	}
	for _, testCase := range testCases {
		req, err := http.NewRequest(
			testCase.Given.Request.Method,
			testCase.Given.Request.Path,
			nil,
		)
		require.NoError(t, err)
		response := rbacServe(t, req,
			testCase.Given.MockResponse.Resource,
			testCase.Given.MockResponse.Verb,
			testCase.Given.MockResponse.Allowed,
			testCase.Given.MockResponse.Err)
		require.NotNil(t, response)
		assert.Equal(t, testCase.Expected.Code, response.Code)
		assert.Equal(t, testCase.Expected.Body, string(response.Body.Bytes()))
	}
}
