package model

import "encoding/json"

var ReportUploadTypes = []ReportUploadType{
	ReportTypeUploadPaymentsMOTOCard,
	ReportTypeUploadPaymentsOnlineCard,
	ReportTypeUploadPaymentsOPGBACS,
	ReportTypeUploadPaymentsSupervisionBACS,
	ReportTypeUploadDebtChase,
	ReportTypeUploadDeputySchedule,
}

type ReportUploadType int

const (
	ReportTypeUploadUnknown ReportUploadType = iota
	ReportTypeUploadPaymentsMOTOCard
	ReportTypeUploadPaymentsOnlineCard
	ReportTypeUploadPaymentsOPGBACS
	ReportTypeUploadPaymentsSupervisionBACS
	ReportTypeUploadDebtChase
	ReportTypeUploadDeputySchedule
)

var reportTypeUploadMap = map[string]ReportUploadType{
	"PaymentsMOTOCard":        ReportTypeUploadPaymentsMOTOCard,
	"PaymentsOnlineCard":      ReportTypeUploadPaymentsOnlineCard,
	"PaymentsOPGBACS":         ReportTypeUploadPaymentsOPGBACS,
	"PaymentsSupervisionBACS": ReportTypeUploadPaymentsSupervisionBACS,
	"DebtChase":               ReportTypeUploadDebtChase,
	"DeputySchedule":          ReportTypeUploadDeputySchedule,
}

func (i ReportUploadType) String() string {
	return i.Key()
}

func (i ReportUploadType) Translation() string {
	switch i {
	case ReportTypeUploadPaymentsMOTOCard:
		return "Payments - MOTO card"
	case ReportTypeUploadPaymentsOnlineCard:
		return "Payments - Online card"
	case ReportTypeUploadPaymentsOPGBACS:
		return "Payments - OPG BACS"
	case ReportTypeUploadPaymentsSupervisionBACS:
		return "Payments - Supervision BACS"
	case ReportTypeUploadDebtChase:
		return "Debt chase"
	case ReportTypeUploadDeputySchedule:
		return "Deputy schedule"
	default:
		return ""
	}
}

func (i ReportUploadType) Key() string {
	switch i {
	case ReportTypeUploadPaymentsMOTOCard:
		return "PaymentsMOTOCard"
	case ReportTypeUploadPaymentsOnlineCard:
		return "PaymentsOnlineCard"
	case ReportTypeUploadPaymentsOPGBACS:
		return "PaymentsOPGBACS"
	case ReportTypeUploadPaymentsSupervisionBACS:
		return "PaymentsSupervisionBACS"
	case ReportTypeUploadDebtChase:
		return "DebtChase"
	case ReportTypeUploadDeputySchedule:
		return "DeputySchedule"
	default:
		return ""
	}
}

func ParseReportUploadType(s string) ReportUploadType {
	value, ok := reportTypeUploadMap[s]
	if !ok {
		return ReportUploadType(0)
	}
	return value
}

func (i ReportUploadType) Valid() bool {
	return i != ReportTypeUploadUnknown
}

func (i ReportUploadType) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.Key())
}

func (i *ReportUploadType) UnmarshalJSON(data []byte) (err error) {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	*i = ParseReportUploadType(s)
	return nil
}
