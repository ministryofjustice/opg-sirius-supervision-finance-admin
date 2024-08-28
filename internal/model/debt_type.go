package model

var ReportDebtTypes = []ReportDebtType{
	ReportDebtTypeFeeChase,
	ReportDebtTypeFinalFee,
}

type ReportDebtType int

const (
	ReportDebtTypeUnknown ReportDebtType = iota
	ReportDebtTypeFeeChase
	ReportDebtTypeFinalFee
)

func (i ReportDebtType) String() string {
	return i.Key()
}

func (i ReportDebtType) Translation() string {
	switch i {
	case ReportDebtTypeFeeChase:
		return "Fee Chase"
	case ReportDebtTypeFinalFee:
		return "Final Fee"
	default:
		return ""
	}
}

func (i ReportDebtType) Key() string {
	switch i {
	case ReportDebtTypeFeeChase:
		return "FeeChase"
	case ReportDebtTypeFinalFee:
		return "FinalFee"
	default:
		return ""
	}
}

func (i ReportDebtType) Valid() bool {
	return i != ReportDebtTypeUnknown
}
