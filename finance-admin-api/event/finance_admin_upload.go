package event

import "context"

type FinanceAdminUpload struct {
	EmailAddress string `json:"emailAddress"`
	Filename     string `json:"filename"`
}

func (c *Client) FinanceAdminUpload(ctx context.Context, event FinanceAdminUpload) error {
	return c.send(ctx, "finance-admin-upload", event)
}
