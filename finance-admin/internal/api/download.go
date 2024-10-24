package api

import (
	"fmt"
	"net/http"
)

func (c *Client) Download(ctx Context, uid string) (*http.Response, error) {
	req, _ := c.newBackendRequest(ctx, http.MethodGet, fmt.Sprintf("/download?uid=%s", uid), nil)

	return c.http.Do(req)
}
