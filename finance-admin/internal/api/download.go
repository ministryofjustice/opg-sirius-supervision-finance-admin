package api

import (
	"fmt"
	"net/http"
)

func (c *Client) Download(ctx Context, uid string) (*http.Response, error) {
	req, err := c.newBackendRequest(ctx, http.MethodGet, fmt.Sprintf("/download?uid=%s", uid), nil)

	if err != nil {
		return nil, err
	}

	return c.http.Do(req)
}
