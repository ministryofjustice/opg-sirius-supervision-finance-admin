package shared

import (
	"encoding/json"
	"fmt"
)

const (
	EventSourceFinanceHub                 = "opg.supervision.finance"
	DetailTypeMotoPaymentsReportProcessed = "moto-payments-report-processed"
)

type Event struct {
	Source       string      `json:"source"`
	EventBusName string      `json:"event-bus-name"`
	DetailType   string      `json:"detail-type"`
	Detail       interface{} `json:"detail"`
}

func (e *Event) UnmarshalJSON(data []byte) error {
	type tmp Event // avoids infinite recursion
	if err := json.Unmarshal(data, (*tmp)(e)); err != nil {
		return err
	}

	var raw struct {
		Detail json.RawMessage `json:"detail"`
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	switch e.DetailType {
	case DetailTypeMotoPaymentsReportProcessed:
		var detail MotoPaymentsReportProcessedEvent
		if err := json.Unmarshal(raw.Detail, &detail); err != nil {
			return err
		}
		e.Detail = detail
	default:
		return fmt.Errorf("unknown detail type: %s", e.DetailType)
	}

	return nil
}

type MotoPaymentsReportProcessedEvent struct {
	ClientID      int    `json:"clientId"`
	CaseRecNumber string `json:"caseRecNumber"`
	EmailAddress  string `json:"emailAddress"`
}

type RequestParameters struct {
	BucketName string `json:"bucketName"`
	Key        string `json:"key"`
}
