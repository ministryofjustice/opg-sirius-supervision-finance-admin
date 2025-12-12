package shared

import "encoding/json"

var DebtTypes = []DebtType{
	DebtTypeFeeChase,
	DebtTypeFinalFee,
	DebtTypeApprovedRefunds,
	DebtTypeAllRefunds,
}

type DebtType int

const (
	DebtTypeUnknown DebtType = iota
	DebtTypeFeeChase
	DebtTypeFinalFee
	DebtTypeApprovedRefunds
	DebtTypeAllRefunds
)

var debtTypeMap = map[string]DebtType{
	"FeeChase":        DebtTypeFeeChase,
	"FinalFee":        DebtTypeFinalFee,
	"ApprovedRefunds": DebtTypeApprovedRefunds,
	"AllRefunds":      DebtTypeAllRefunds,
}

func (d DebtType) String() string {
	return d.Key()
}

func (d DebtType) Translation() string {
	switch d {
	case DebtTypeFeeChase:
		return "Fee Chase"
	case DebtTypeFinalFee:
		return "Final Fee Debt"
	case DebtTypeApprovedRefunds:
		return "Billing only - Approved refunds"
	case DebtTypeAllRefunds:
		return "All Refunds"
	default:
		return ""
	}
}

func (d DebtType) Key() string {
	switch d {
	case DebtTypeFeeChase:
		return "FeeChase"
	case DebtTypeFinalFee:
		return "FinalFee"
	case DebtTypeApprovedRefunds:
		return "ApprovedRefunds"
	case DebtTypeAllRefunds:
		return "AllRefunds"
	default:
		return ""
	}
}

func ParseReportDebtType(s string) *DebtType {
	value, ok := debtTypeMap[s]
	if !ok {
		return nil
	}
	return &value
}

func (d DebtType) Valid() bool {
	return d != DebtTypeUnknown
}

func (d DebtType) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.Key())
}
