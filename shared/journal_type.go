package shared

import "encoding/json"

var JournalTypes = []JournalType{
	JournalTypeReceiptTransactions,
	JournalTypeNonReceiptTransactions,
	JournalTypeUnappliedTransactions,
}

const (
	JournalTypeUnknown JournalType = iota
	JournalTypeReceiptTransactions
	JournalTypeNonReceiptTransactions
	JournalTypeUnappliedTransactions
)

var journalTypeMap = map[string]JournalType{
	"ReceiptTransactions":    JournalTypeReceiptTransactions,
	"NonReceiptTransactions": JournalTypeNonReceiptTransactions,
	"UnappliedTransactions":  JournalTypeUnappliedTransactions,
}

type JournalType int

func (j JournalType) String() string {
	return j.Key()
}

func (j JournalType) Translation() string {
	switch j {
	case JournalTypeReceiptTransactions:
		return "Receipt Transactions"
	case JournalTypeNonReceiptTransactions:
		return "Non Receipt Transactions"
	case JournalTypeUnappliedTransactions:
		return "Refunds & Unapplied Transactions"
	default:
		return ""
	}
}

func (j JournalType) Key() string {
	switch j {
	case JournalTypeReceiptTransactions:
		return "ReceiptTransactions"
	case JournalTypeNonReceiptTransactions:
		return "NonReceiptTransactions"
	case JournalTypeUnappliedTransactions:
		return "UnappliedTransactions"
	default:
		return ""
	}
}

func ParseJournalType(s string) *JournalType {
	value, ok := journalTypeMap[s]
	if !ok {
		return nil
	}
	return &value
}

func (j JournalType) Valid() bool {
	return j != JournalTypeUnknown
}

func (j JournalType) MarshalJSON() ([]byte, error) {
	return json.Marshal(j.Key())
}
