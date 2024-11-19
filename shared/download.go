package shared

import (
	"encoding/base64"
	"encoding/json"
)

type Download struct {
	Uid               string `json:"uid"`
	ReportType        string `json:"reportType"`
	ReportAccountType string `json:"reportAccountType"`
	DateOfTransaction string `json:"dateOfTransaction"`
}

type DownloadRequest struct {
	Key       string
	VersionId string
}

func (d *DownloadRequest) Encode() (string, error) {
	jsonData, err := json.Marshal(d)
	if err != nil {
		return "", err
	}

	base64Data := base64.StdEncoding.EncodeToString(jsonData)
	return base64Data, nil
}

func (d *DownloadRequest) Decode(data string) error {
	jsonData, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return err
	}

	err = json.Unmarshal(jsonData, &d)
	if err != nil {
		return err
	}

	return nil
}
