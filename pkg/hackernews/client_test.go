package hackernews_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/emmaLP/gs-software-onboarding/pkg/hackernews"
	"github.com/emmaLP/gs-software-onboarding/pkg/hackernews/model"
	"github.com/stretchr/testify/assert"
)

type testConfig struct {
	responseBody []byte
	statusCode   int
	expectedErr  string
}

func TestGetTopStories(t *testing.T) {
	successfulResponseBody, _ := json.Marshal([]int{1, 2, 3, 4})
	blankResponseBody, _ := json.Marshal([]int{})

	tests := map[string]testConfig{
		"Successfully get story ids": {
			responseBody: successfulResponseBody,
			statusCode:   http.StatusOK,
		},
		"Request Failed with InternalServerError": {
			responseBody: blankResponseBody,
			statusCode:   http.StatusInternalServerError,
			expectedErr:  "Unexpected status code returned. Got 500 status code",
		},
		"Invalid json body": {
			responseBody: []byte("Hello"),
			statusCode:   http.StatusOK,
			expectedErr:  "Unabled to decode response: invalid character 'H' looking for beginning of value",
		},
	}
	for testName, testValues := range tests {
		t.Run(testName, func(t *testing.T) {
			server := testServer(testValues.statusCode, testValues.responseBody, t)
			defer server.Close()
			client, err := hackernews.New(server.URL, server.Client())
			assert.Equal(t, nil, err, "Failed to create hackernews client")

			result, err := client.GetTopStories()
			if testValues.expectedErr != "" {
				assert.EqualErrorf(t, err, testValues.expectedErr, "Request failed should be: %v, got: %v", testValues.expectedErr, err)
				assert.Nil(t, result)
			} else {
				assert.Nil(t, err)
				var expectedResult []int
				_ = json.Unmarshal(testValues.responseBody, &expectedResult)
				assert.Equal(t, expectedResult, result)
			}
		})
	}
}

func TestGetItem(t *testing.T) {
	successfulResponseBody, _ := json.Marshal(model.Item{
		ID:        1,
		Type:      "story",
		Text:      "<i>or</i> HN: the Next Iteration<p>I get the impression that with Arc being released a lot of people who never had time for HN before are suddenly dropping in more often. (PG: what are the numbers on this? I'm envisioning a spike.)<p>Not to say that isn't great, but I'm wary of Diggification. Between links comparing programming to sex and a flurry of gratuitous, ostentatious  adjectives in the headlines it's a bit concerning.<p>80% of the stuff that makes the front page is still pretty awesome, but what's in place to keep the signal/noise ratio high? Does the HN model still work as the community scales? What's in store for (++ HN)?",
		URL:       "",
		Score:     25,
		Title:     "Ask HN: The Arc Effect",
		CreatedAt: time.Time{},
		CreatedBy: "tel",
		Dead:      false,
		Deleted:   false,
	})
	blankResponseBody, _ := json.Marshal(model.Item{})

	tests := map[string]testConfig{
		"Successfully get story ids": {
			responseBody: successfulResponseBody,
			statusCode:   http.StatusOK,
		},
		"Request Failed with InternalServerError": {
			responseBody: blankResponseBody,
			statusCode:   http.StatusInternalServerError,
			expectedErr:  "Unexpected status code returned. Got 500 status code",
		},
		"Invalid json body": {
			responseBody: []byte("Hello"),
			statusCode:   http.StatusOK,
			expectedErr:  "Unabled to decode response: invalid character 'H' looking for beginning of value",
		},
	}
	counter := 1
	for testName, testValues := range tests {
		t.Run(testName, func(t *testing.T) {
			server := testServer(testValues.statusCode, testValues.responseBody, t)
			defer server.Close()
			client, err := hackernews.New(server.URL, server.Client())
			assert.Equal(t, nil, err, "Failed to create hackernews client")

			result, err := client.GetItem(counter)
			if testValues.expectedErr != "" {
				assert.EqualErrorf(t, err, testValues.expectedErr, "Request failed should be: %v, got: %v", testValues.expectedErr, err)
				assert.Nil(t, result)
			} else {
				var expectedResult *model.Item
				_ = json.Unmarshal(testValues.responseBody, &expectedResult)
				assert.Equal(t, expectedResult, result)
			}
		})
		counter++
	}
}

func testServer(statusCode int, responseBody []byte, t *testing.T) *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Send response to be tested
		rw.WriteHeader(statusCode)
		_, err := rw.Write(responseBody)
		assert.Equal(t, nil, err, "Failed to write the body")
	}))

	return server
}
