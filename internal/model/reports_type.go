package model

var ReportsTypes = []ReportsType{
	ReportsTypeJournal,
	ReportsTypeSchedule,
	ReportsTypeAccountsReceivable,
	ReportsTypeDebt,
}

type ReportsType int

const (
	ReportsTypeUnknown ReportsType = iota
	ReportsTypeJournal
	ReportsTypeSchedule
	ReportsTypeAccountsReceivable
	ReportsTypeDebt
)

func (i ReportsType) String() string {
	return i.Key()
}

func (i ReportsType) Translation() string {
	switch i {
	case ReportsTypeJournal:
		return "Journal"
	case ReportsTypeSchedule:
		return "Schedule"
	case ReportsTypeAccountsReceivable:
		return "Accounts Receivable"
	case ReportsTypeDebt:
		return "Debt"
	default:
		return ""
	}
}

func (i ReportsType) Key() string {
	switch i {
	case ReportsTypeJournal:
		return "Journal"
	case ReportsTypeSchedule:
		return "Schedule"
	case ReportsTypeAccountsReceivable:
		return "AccountsReceivable"
	case ReportsTypeDebt:
		return "Debt"
	default:
		return ""
	}
}

func (i ReportsType) Valid() bool {
	return i != ReportsTypeUnknown
}

func (i ReportsType) RequiresDateValidation() bool {
	switch i {
	case ReportsTypeJournal:
		return true
	}
	return false
}
