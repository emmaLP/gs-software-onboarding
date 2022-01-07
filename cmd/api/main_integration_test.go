//go:build integration

package main

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	commonModel "github.com/emmaLP/gs-software-onboarding/pkg/common/model"
	"github.com/emmaLP/gs-software-onboarding/pkg/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListAll(t *testing.T) {
	handler := test.LoadTestHandler(t, context.TODO())

	handler.SaveItemToDatabase(context.TODO(), &commonModel.Item{
		ID:      2,
		Dead:    true,
		Deleted: false,
		Type:    "job",
	})
	handler.SaveItemToDatabase(context.TODO(), &commonModel.Item{
		ID:      1,
		Dead:    false,
		Deleted: false,
		Type:    "story",
	})
	httpHelper := &test.Http{
		HTTPClient: http.DefaultClient,
	}

	body := &test.HttpResponse{Items: []commonModel.Item{}}

	err := httpHelper.GetRequest(fmt.Sprintf("http://localhost%s/all", handler.Config.Api.Address), &body)
	require.NoError(t, err)
	assert.Len(t, body.Items, 2)
}
