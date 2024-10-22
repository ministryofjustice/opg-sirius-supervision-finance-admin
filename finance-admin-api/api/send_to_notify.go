package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/opg-sirius-finance-admin/shared"
	"net/http"
	"os"
	"slices"
	"strings"
	"time"
)

const notifyUrl = "https://api.notifications.service.gov.uk"
const emailEndpoint = "v2/notifications/email"
const processingErrorTemplateId = "872d88b3-076e-495c-bf81-a2be2d3d234c"
const processingFailedTemplateId = "a8f9ab79-1489-4639-9e6c-cad1f079ebcf"
const processingSuccessTemplateId = "8c85cf6c-695f-493a-a25f-77b4fb5f6a8e"

type ProcessingFailedPersonalisation struct {
	FailedLines []string `json:"failed_lines"`
	ReportType  string   `json:"report_type"`
}

type ProcessingSuccessPersonalisation struct {
	ReportType string `json:"report_type"`
}

type NotifyPayload struct {
	EmailAddress    string      `json:"email_address"`
	TemplateId      string      `json:"template_id"`
	Personalisation interface{} `json:"personalisation"`
}

func parseNotifyApiKey(notifyApiKey string) (string, string) {
	splitKey := strings.Split(notifyApiKey, "-")
	if len(splitKey) != 11 {
		return "", ""
	}
	iss := fmt.Sprintf("%s-%s-%s-%s-%s", splitKey[1], splitKey[2], splitKey[3], splitKey[4], splitKey[5])
	jwtToken := fmt.Sprintf("%s-%s-%s-%s-%s", splitKey[6], splitKey[7], splitKey[8], splitKey[9], splitKey[10])
	return iss, jwtToken
}

func createSignedJwtToken() (string, error) {
	iss, jwtKey := parseNotifyApiKey(os.Getenv("OPG_NOTIFY_API_KEY"))

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": iss,
		"iat": time.Now().Unix(),
	})

	signedToken, err := t.SignedString([]byte(jwtKey))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func formatFailedLines(failedLines map[int]string) []string {
	var errorMessage string
	var formattedLines []string
	var keys []int
	for i := range failedLines {
		keys = append(keys, i)
	}

	slices.Sort(keys)

	for _, key := range keys {
		failedLine := failedLines[key]
		errorMessage = ""

		switch failedLine {
		case "DATE_PARSE_ERROR":
			errorMessage = "Unable to parse date"
		case "DUPLICATE_PAYMENT":
			errorMessage = "Duplicate payment line"
		case "CLIENT_NOT_FOUND":
			errorMessage = "Could not find a client with this court reference"
		}

		formattedLines = append(formattedLines, fmt.Sprintf("Line %d: %s", key, errorMessage))
	}

	return formattedLines
}

func createNotifyPayload(detail shared.FinanceAdminUploadProcessedEvent, reportType string) NotifyPayload {
	var payload NotifyPayload

	if detail.Error != "" {
		payload = NotifyPayload{
			detail.EmailAddress,
			processingErrorTemplateId,
			struct {
				Error      string `json:"error"`
				ReportType string `json:"report_type"`
			}{
				detail.Error,
				reportType,
			},
		}
	} else if len(detail.FailedLines) != 0 {
		payload = NotifyPayload{
			detail.EmailAddress,
			processingFailedTemplateId,
			struct {
				FailedLines []string `json:"failed_lines"`
				ReportType  string   `json:"report_type"`
			}{
				formatFailedLines(detail.FailedLines),
				reportType,
			},
		}
	} else {
		payload = NotifyPayload{
			detail.EmailAddress,
			processingSuccessTemplateId,
			struct {
				ReportType string `json:"report_type"`
			}{reportType},
		}
	}

	return payload
}

func (s *Server) SendEmailToNotify(ctx context.Context, detail shared.FinanceAdminUploadProcessedEvent, reportType string) error {
	signedToken, err := createSignedJwtToken()
	if err != nil {
		return err
	}

	var body bytes.Buffer

	err = json.NewEncoder(&body).Encode(createNotifyPayload(detail, reportType))
	if err != nil {
		return err
	}

	r, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/%s", notifyUrl, emailEndpoint), &body)

	if err != nil {
		return err
	}

	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Authorization", "Bearer "+signedToken)

	resp, err := s.http.Do(r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusCreated {
		return nil
	}

	return newStatusError(resp)
}
