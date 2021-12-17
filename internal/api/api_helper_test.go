package api

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	hnModel "github.com/emmaLP/gs-software-onboarding/pkg/hackernews/model"
	"github.com/labstack/echo/v4"
)

// Reusable set of functions that to be used in api related tests

func setupRequest(t *testing.T, path string) (*httptest.ResponseRecorder, echo.Context) {
	t.Helper()
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath(path)

	return rec, c
}

type successResponse struct {
	Items []hnModel.Item `json:"items"`
}

func decodeRequest(t *testing.T, body io.Reader) (res successResponse) {
	t.Helper()

	err := json.NewDecoder(body).Decode(&res)
	if err != nil {
		t.Fatalf("unable to decode response body. body: %v, error: %v", body, err)
	}

	return res
}
