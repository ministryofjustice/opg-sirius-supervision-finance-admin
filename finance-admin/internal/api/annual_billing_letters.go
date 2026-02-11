package api

import (
	"context"
	"fmt"
	"net/http"
)

func (c *Client) AnnualBillingLetters(ctx context.Context) (*http.Response, error) {
	req, err := c.newHubRequest(ctx, http.MethodGet, fmt.Sprintf("/annual-billing-letters"), nil)

	if err != nil {
		return nil, err
	}
	return c.http.Do(req)
}
