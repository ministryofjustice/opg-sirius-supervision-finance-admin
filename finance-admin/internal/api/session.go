package api

import (
	"net/http"
)

func (c *Client) CheckUserSession(ctx Context) (bool, error) {
	req, _ := c.newSessionRequest(ctx)

	res, err := c.http.Do(req)

	return err == nil && res.StatusCode == http.StatusOK, err
}
