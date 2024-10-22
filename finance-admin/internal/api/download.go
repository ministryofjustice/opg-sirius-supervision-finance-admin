package api

import (
	"fmt"
	"net/http"
)

func (c *Client) Download(ctx Context, filename string) (*http.Response, error) {
	req, _ := c.newBackendRequest(ctx, http.MethodGet, fmt.Sprintf("/downloads/%s", filename), nil)

	return c.http.Do(req)
}
