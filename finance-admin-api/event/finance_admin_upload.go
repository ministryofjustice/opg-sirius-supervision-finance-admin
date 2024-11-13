package event

import (
	"context"
	"github.com/opg-sirius-finance-admin/shared"
)

type FinanceAdminUpload struct {
	EmailAddress string      `json:"emailAddress"`
	Filename     string      `json:"filename"`
	UploadType   string      `json:"uploadType"`
	UploadDate   shared.Date `json:"uploadDate"`
}

func (c *Client) FinanceAdminUpload(ctx context.Context, event FinanceAdminUpload) error {
	return c.send(ctx, "finance-admin-upload", event)
}
