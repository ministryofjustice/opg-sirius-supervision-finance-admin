package api

import (
	"context"
	"net/http"
)

func (c *Client) AnnualBillingLetters(ctx context.Context) (*http.Response, error) {
	req, err := c.newHubRequest(ctx, http.MethodGet, "/annual-billing-letters", nil)

	if err != nil {
		return nil, err
	}
	return c.http.Do(req)
}
