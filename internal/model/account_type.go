package model

var ReportAccountTypes = []ReportAccountType{
	ReportAccountTypeAgedDebt,
	ReportAccountTypeUnappliedReceipts,
	ReportAccountTypeCustomerAgeingBuckets,
	ReportAccountTypeARPaidInvoiceReport,
	ReportAccountTypePaidInvoiceTransactionLines,
	ReportAccountTypeTotalReceiptsReport,
	ReportAccountTypeBadDebtWriteOffReport,
	ReportAccountTypeFeeAccrual,
}

type ReportAccountType int

const (
	ReportAccountTypeUnknown ReportAccountType = iota
	ReportAccountTypeAgedDebt
	ReportAccountTypeUnappliedReceipts
	ReportAccountTypeCustomerAgeingBuckets
	ReportAccountTypeARPaidInvoiceReport
	ReportAccountTypePaidInvoiceTransactionLines
	ReportAccountTypeTotalReceiptsReport
	ReportAccountTypeBadDebtWriteOffReport
	ReportAccountTypeFeeAccrual
)

func (i ReportAccountType) String() string {
	return i.Key()
}

func (i ReportAccountType) Translation() string {
	switch i {
	case ReportAccountTypeAgedDebt:
		return "Aged Debt"
	case ReportAccountTypeUnappliedReceipts:
		return "Unapplied Receipts"
	case ReportAccountTypeCustomerAgeingBuckets:
		return "Customer Ageing Buckets"
	case ReportAccountTypeARPaidInvoiceReport:
		return "AR Paid Invoice Report"
	case ReportAccountTypePaidInvoiceTransactionLines:
		return "Paid Invoice Transaction Lines"
	case ReportAccountTypeTotalReceiptsReport:
		return "Total Receipts Report"
	case ReportAccountTypeBadDebtWriteOffReport:
		return "Bad Debt Write-off Report"
	case ReportAccountTypeFeeAccrual:
		return "Fee Accrual"
	default:
		return ""
	}
}

func (i ReportAccountType) Key() string {
	switch i {
	case ReportAccountTypeAgedDebt:
		return "AgedDebt"
	case ReportAccountTypeUnappliedReceipts:
		return "UnappliedReceipts"
	case ReportAccountTypeCustomerAgeingBuckets:
		return "CustomerAgeingBuckets"
	case ReportAccountTypeARPaidInvoiceReport:
		return "ARPaidInvoiceReport"
	case ReportAccountTypePaidInvoiceTransactionLines:
		return "PaidInvoiceTransactionLines"
	case ReportAccountTypeTotalReceiptsReport:
		return "TotalReceiptsReport"
	case ReportAccountTypeBadDebtWriteOffReport:
		return "BadDebtWriteOffReport"
	case ReportAccountTypeFeeAccrual:
		return "FeeAccrual"
	default:
		return ""
	}
}

func (i ReportAccountType) Valid() bool {
	return i != ReportAccountTypeUnknown
}
