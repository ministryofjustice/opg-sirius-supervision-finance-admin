package model

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

func (i ReportUploadType) Valid() bool {
	return i != ReportTypeUploadUnknown
}
