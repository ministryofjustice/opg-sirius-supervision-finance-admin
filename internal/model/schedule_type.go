package model

var ReportScheduleTypes = []ReportScheduleType{
	ReportTypeOnlineCardPayments,
	ReportOPGBACSTransfer,
	ReportSupervisionBACSTransfer,
	ReportDirectDebitPayments,
	ReportAdFeeInvoices,
	ReportS2FeeInvoices,
	ReportS3FeeInvoices,
	ReportB2FeeInvoices,
	ReportB3FeeInvoices,
	ReportSFFeeInvoicesGeneral,
	ReportSFFeeInvoicesMinimal,
	ReportSEFeeInvoicesGeneral,
	ReportSEFeeInvoicesMinimal,
	ReportSOFeeInvoicesGeneral,
	ReportSOFeeInvoicesMinimal,
	ReportADFeeReductions,
	ReportGeneralManualCredits,
	ReportMinimalManualCredits,
	ReportGeneralManualDebits,
	ReportMinimalManualDebits,
	ReportADWriteOffs,
	ReportGeneralWriteOffs,
	ReportMinimalWriteOffs,
	ReportADWriteOffReversals,
	ReportGeneralWriteOffReversals,
	ReportMinimalWriteOffReversals,
}

type ReportScheduleType int

const (
	ReportScheduleTypeUnknown ReportScheduleType = iota
	ReportTypeOnlineCardPayments
	ReportOPGBACSTransfer
	ReportSupervisionBACSTransfer
	ReportDirectDebitPayments
	ReportAdFeeInvoices
	ReportS2FeeInvoices
	ReportS3FeeInvoices
	ReportB2FeeInvoices
	ReportB3FeeInvoices
	ReportSFFeeInvoicesGeneral
	ReportSFFeeInvoicesMinimal
	ReportSEFeeInvoicesGeneral
	ReportSEFeeInvoicesMinimal
	ReportSOFeeInvoicesGeneral
	ReportSOFeeInvoicesMinimal
	ReportADFeeReductions
	ReportGeneralManualCredits
	ReportMinimalManualCredits
	ReportGeneralManualDebits
	ReportMinimalManualDebits
	ReportADWriteOffs
	ReportGeneralWriteOffs
	ReportMinimalWriteOffs
	ReportADWriteOffReversals
	ReportGeneralWriteOffReversals
	ReportMinimalWriteOffReversals
)

func (i ReportScheduleType) String() string {
	return i.Key()
}

func (i ReportScheduleType) Translation() string {
	switch i {
	case ReportTypeOnlineCardPayments:
		return "Online Card Payments"
	case ReportOPGBACSTransfer:
		return "OPG BACS Transfer"
	case ReportSupervisionBACSTransfer:
		return "Supervision BACS transfer"
	case ReportDirectDebitPayments:
		return "Direct Debit Payment"
	case ReportAdFeeInvoices:
		return "Ad Fee Invoices"
	case ReportS2FeeInvoices:
		return "S2 Fee Invoices"
	case ReportS3FeeInvoices:
		return "S3 Fee Invoices"
	case ReportB2FeeInvoices:
		return "B2 Fee Invoices"
	case ReportB3FeeInvoices:
		return "B3 Fee Invoices"
	case ReportSFFeeInvoicesGeneral:
		return "SF Fee Invoices (General) "
	case ReportSFFeeInvoicesMinimal:
		return "SF Fee Invoices (Minimal)"
	case ReportSEFeeInvoicesGeneral:
		return "SE Fee Invoices (General)"
	case ReportSEFeeInvoicesMinimal:
		return "SE Fee Invoices (Minimal)"
	case ReportSOFeeInvoicesGeneral:
		return "SO Fee Invoices (General)"
	case ReportSOFeeInvoicesMinimal:
		return "SO Fee Invoices (Minimal)"
	case ReportADFeeReductions:
		return "AD Fee Reductions"
	case ReportGeneralManualCredits:
		return "General Manual Credits"
	case ReportMinimalManualCredits:
		return "Minimal Manual Credits"
	case ReportGeneralManualDebits:
		return "General Manual Debits"
	case ReportMinimalManualDebits:
		return "Minimal Manual Debits"
	case ReportADWriteOffs:
		return "AD Write-offs"
	case ReportGeneralWriteOffs:
		return "General Write-offs"
	case ReportMinimalWriteOffs:
		return "Minimal Write-offs"
	case ReportADWriteOffReversals:
		return "AD Write-off Reversals"
	case ReportGeneralWriteOffReversals:
		return "General Write-off Reversals"
	case ReportMinimalWriteOffReversals:
		return "Minimal Write-off Reversals"
	default:
		return ""
	}
}

func (i ReportScheduleType) Key() string {
	switch i {
	case ReportTypeOnlineCardPayments:
		return "OnlineCardPayments"
	case ReportOPGBACSTransfer:
		return "OPGBACSTransfer"
	case ReportSupervisionBACSTransfer:
		return "SupervisionBACSTransfer"
	case ReportDirectDebitPayments:
		return "DirectDebitPayment"
	case ReportAdFeeInvoices:
		return "AdFeeInvoices"
	case ReportS2FeeInvoices:
		return "S2FeeInvoices"
	case ReportS3FeeInvoices:
		return "S3FeeInvoices"
	case ReportB2FeeInvoices:
		return "B2FeeInvoices"
	case ReportB3FeeInvoices:
		return "B3FeeInvoices"
	case ReportSFFeeInvoicesGeneral:
		return "SFFeeInvoicesGeneral "
	case ReportSFFeeInvoicesMinimal:
		return "SFFeeInvoicesMinimal"
	case ReportSEFeeInvoicesGeneral:
		return "SEFeeInvoicesGeneral"
	case ReportSEFeeInvoicesMinimal:
		return "SEFeeInvoicesMinimal"
	case ReportSOFeeInvoicesGeneral:
		return "SOFeeInvoicesGeneral"
	case ReportSOFeeInvoicesMinimal:
		return "SOFeeInvoicesMinimal"
	case ReportADFeeReductions:
		return "ADFeeReductions"
	case ReportGeneralManualCredits:
		return "GeneralManualCredits"
	case ReportMinimalManualCredits:
		return "MinimalManualCredits"
	case ReportGeneralManualDebits:
		return "GeneralManualDebits"
	case ReportMinimalManualDebits:
		return "MinimalManualDebits"
	case ReportADWriteOffs:
		return "ADWrite-offs"
	case ReportGeneralWriteOffs:
		return "GeneralWrite-offs"
	case ReportMinimalWriteOffs:
		return "MinimalWriteOffs"
	case ReportADWriteOffReversals:
		return "ADWriteOffReversals"
	case ReportGeneralWriteOffReversals:
		return "GeneralWriteOffReversals"
	case ReportMinimalWriteOffReversals:
		return "MinimalWriteOffReversals"
	default:
		return ""
	}
}

func (i ReportScheduleType) Valid() bool {
	return i != ReportScheduleTypeUnknown
}
