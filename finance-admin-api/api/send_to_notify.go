package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

const notifyUrl = "https://api.notifications.service.gov.uk"
const emailEndpoint = "v2/notifications/email"

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

func (s *Server) SendEmailToNotify(ctx context.Context, emailAddress string, templateId string, failedLines []string, reportType string) error {
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
		Personalisation: Personalisation{failedLines, reportType},
	}

	fmt.Println(payload)

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

	fmt.Println(resp.Body)

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodyString := string(bodyBytes)
	fmt.Println(bodyString)
	_, _ = io.Copy(os.Stdout, resp.Body)
	return nil
}
