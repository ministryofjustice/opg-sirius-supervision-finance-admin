package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"os"
	"strings"
	"time"
)

const notifyUrl = "https://api.notifications.service.gov.uk"
const emailEndpoint = "v2/notifications/email"
const templateId = "a8f9ab79-1489-4639-9e6c-cad1f079ebcf"

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

	for i, line := range failedLines {
		errorMessage = ""

		switch line {
		case "DATE_PARSE_ERROR":
			errorMessage = "Unable to parse date"
		case "DUPLICATE_PAYMENT":
			errorMessage = "Duplicate payment line"
		case "CLIENT_NOT_FOUND":
			errorMessage = "Could not find a client with this court reference"
		}

		formattedLines = append(formattedLines, fmt.Sprintf("Line %d: %s", i, errorMessage))
	}

	return formattedLines
}

func (s *Server) SendEmailToNotify(ctx context.Context, emailAddress string, failedLines map[int]string, reportType string) error {
	signedToken, err := createSignedJwtToken()
	if err != nil {
		return err
	}

	type Personalisation struct {
		FailedLines []string `json:"failed_lines"`
		ReportType  string   `json:"report_type"`
	}

	payload := struct {
		EmailAddress    string          `json:"email_address"`
		TemplateId      string          `json:"template_id"`
		Personalisation Personalisation `json:"personalisation"`
	}{
		EmailAddress:    emailAddress,
		TemplateId:      templateId,
		Personalisation: Personalisation{formatFailedLines(failedLines), reportType},
	}

	var body bytes.Buffer

	err = json.NewEncoder(&body).Encode(payload)
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
