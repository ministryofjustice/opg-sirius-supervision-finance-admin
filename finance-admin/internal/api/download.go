package api

import (
	"context"
	"fmt"
	"net/http"
)

func (c *Client) Download(ctx context.Context, uid string) (*http.Response, error) {
	req, err := c.newHubRequest(ctx, http.MethodGet, fmt.Sprintf("/download?uid=%s", uid), nil)

	if err != nil {
		return nil, err
	}

	return c.http.Do(req)
}
