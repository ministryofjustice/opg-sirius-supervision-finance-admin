package model

var ReportJournalTypes = []ReportJournalType{
	ReportTypeReceiptTransactions,
	ReportTypeNonReceiptTransactions,
	ReportTypeMOTOCardPayments,
}

type ReportJournalType int

const (
	ReportTypeUnknown ReportJournalType = iota
	ReportTypeReceiptTransactions
	ReportTypeNonReceiptTransactions
	ReportTypeMOTOCardPayments
)

func (i ReportJournalType) String() string {
	return i.Key()
}

func (i ReportJournalType) Translation() string {
	switch i {
	case ReportTypeReceiptTransactions:
		return "Receipt Transactions"
	case ReportTypeNonReceiptTransactions:
		return "Non Receipt Transactions"
	case ReportTypeMOTOCardPayments:
		return "Accounts Receivable"
	default:
		return ""
	}
}

func (i ReportJournalType) Key() string {
	switch i {
	case ReportTypeReceiptTransactions:
		return "ReceiptTransactions"
	case ReportTypeNonReceiptTransactions:
		return "NonReceiptTransactions"
	case ReportTypeMOTOCardPayments:
		return "AccountsReceivable"
	default:
		return ""
	}
}

func (i ReportJournalType) Valid() bool {
	return i != ReportTypeUnknown
}
