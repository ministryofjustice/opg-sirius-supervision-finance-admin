package api

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func (c *Client) Download(ctx Context, uid string) (*http.Response, error) {
	var body bytes.Buffer
	err := json.NewEncoder(&body).Encode(uid)
	if err != nil {
		return nil, err
	}

	req, err := c.newBackendRequest(ctx, http.MethodPost, "/downloads", &body)

	if err != nil {
		return nil, err
	}

	return c.http.Do(req)
}
