package shared

import "encoding/json"

var JournalTypes = []JournalType{
	JournalTypeReceiptTransactions,
	JournalTypeNonReceiptTransactions,
	JournalTypeNonReceiptTransactionsHistoric,
	JournalTypeReceiptTransactionsHistoric,
}

const (
	JournalTypeUnknown JournalType = iota
	JournalTypeReceiptTransactions
	JournalTypeNonReceiptTransactions
	JournalTypeNonReceiptTransactionsHistoric
	JournalTypeReceiptTransactionsHistoric
)

var journalTypeMap = map[string]JournalType{
	"ReceiptTransactions":            JournalTypeReceiptTransactions,
	"NonReceiptTransactions":         JournalTypeNonReceiptTransactions,
	"NonReceiptTransactionsHistoric": JournalTypeNonReceiptTransactionsHistoric,
	"ReceiptTransactionsHistoric":    JournalTypeReceiptTransactionsHistoric,
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
	case JournalTypeNonReceiptTransactionsHistoric:
		return "Non Receipt Transactions (Historic)"
	case JournalTypeReceiptTransactionsHistoric:
		return "Receipt Transactions (Historic)"
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
	case JournalTypeNonReceiptTransactionsHistoric:
		return "NonReceiptTransactionsHistoric"
	case JournalTypeReceiptTransactionsHistoric:
		return "ReceiptTransactionsHistoric"
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
