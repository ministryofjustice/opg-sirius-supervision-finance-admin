package api

import (
	"net/http"
)

func (c *Client) CheckUserSession(ctx Context) (*http.Response, error, bool) {
	req, _ := c.newSessionRequest(ctx)

	res, err := c.http.Do(req)

	return res, err, err == nil && res.StatusCode == http.StatusOK
}
