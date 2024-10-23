package api

import (
	"context"
	"net/http"
)

func (c *Client) CheckUserSession(ctx context.Context, sessionCookie *http.Cookie) (*http.Response, error, bool) {
	req, _ := c.newSessionRequest(ctx, sessionCookie)

	res, err := c.http.Do(req)

	return res, err, err == nil && res.StatusCode == http.StatusOK
}
