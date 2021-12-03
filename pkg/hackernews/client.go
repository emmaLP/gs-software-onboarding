package hackernews

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/emmaLP/gs-software-onboarding/pkg/hackernews/model"
)

type Client interface {
	GetTopStories() ([]int, error)
	GetItem(id int) (*model.Item, error)
}

type client struct {
	httpClient *http.Client
	baseUrl    string
}

const (
	topStoriesPath = "%s/topstories.json"
	itemPath       = "%s/item/%d.json"
)

func New(baseUrl string, c *http.Client) (*client, error) {
	if c == nil {
		c = http.DefaultClient
	}

	if len(baseUrl) == 0 {
		return nil, errors.New("baseURL cannot be nil or empty")
	}

	return &client{httpClient: c, baseUrl: baseUrl}, nil
}

func (c *client) GetTopStories() ([]int, error) {
	path := fmt.Sprintf(topStoriesPath, c.baseUrl)
	var ids []int
	if err := c.performRequest(path, http.MethodGet, &ids); err != nil {
		return nil, err
	}

	return ids, nil
}
func (c *client) GetItem(id int) (*model.Item, error) {
	path := fmt.Sprintf(itemPath, c.baseUrl, id)
	var item *model.Item

	if err := c.performRequest(path, http.MethodGet, &item); err != nil {
		return nil, err
	}

	return item, nil
}

func (c *client) performRequest(path, method string, result interface{}) error {
	request, err := http.NewRequest(method, path, nil)
	if err != nil {
		return fmt.Errorf("Failed to create http request object: %w", err)
	}

	response, err := c.httpClient.Do(request)
	if err != nil {
		return fmt.Errorf("Unable to make http call: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("Unexpected status code returned. Got %d status code", response.StatusCode)
	}

	if err := json.NewDecoder(response.Body).Decode(result); err != nil {
		return fmt.Errorf("Unabled to decode response: %w", err)
	}
	return nil
}
