package api

import (
	"fmt"
	"github.com/opg-sirius-finance-admin/apierror"
	"net/http"
)

func (c *Client) CheckDownload(ctx Context, uid string) error {
	req, err := c.newBackendRequest(ctx, http.MethodHead, fmt.Sprintf("/download?uid=%s", uid), nil)

	if err != nil {
		return err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}

	switch resp.StatusCode {
	case http.StatusOK:
		return nil
	case http.StatusNotFound:
		return apierror.NotFound{}
	default:
		return newStatusError(resp)
	}
}