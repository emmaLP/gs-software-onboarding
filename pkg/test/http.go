//go:build integration

package test

import (
	"encoding/json"
	"fmt"
	"net/http"

	commonModel "github.com/emmaLP/gs-software-onboarding/pkg/common/model"
)

type HttpResponse struct {
	Items        []commonModel.Item `json:"items,omitempty"`
	Err          error              `json:"error,omitempty"`
	ErrorMessage string             `json:"error_message,omitempty"`
}

type Http struct {
	HTTPClient *http.Client
}

func (c *Http) GetRequest(urlPath string, result interface{}) error {
	if err := c.performRequest(urlPath, http.MethodGet, &result); err != nil {
		return err
	}
	return nil
}

func (c *Http) performRequest(path, method string, result *interface{}) error {
	request, err := http.NewRequest(method, path, nil)
	if err != nil {
		return fmt.Errorf("Failed to create http request object: %w", err)
	}

	response, err := c.HTTPClient.Do(request)
	if err != nil {
		return fmt.Errorf("Unable to make http call: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("Unexpected status code returned. Got %d status code", response.StatusCode)
	}

	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return fmt.Errorf("Unabled to decode response: %w", err)
	}
	return nil
}
