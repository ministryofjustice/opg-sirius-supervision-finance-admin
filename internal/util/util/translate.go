package util

import "github.com/opg-sirius-finance-admin/internal/model"

type pair struct {
	k string
	v string
}

var validationMappings = map[string]map[string]pair{
	"ReportType": {
		"required": pair{"ReportType", "Please select a report type"},
	},
	"ReportSubType": {
		"required": pair{"ReportSubType", "Please select a report to download"},
	},
	"Date": {
		"Date":             pair{"Date", "Please select the report date"},
		"date-in-the-past": pair{"Date", "The report date must be today or in the past"},
	},
	"FromDate": {
		"FromDate": pair{"FromDate", "Date From must be before Date To"},
	},
	"ToDate": {
		"ToDate": pair{"ToDate", "Date To must be after Date From"},
	},
}

func RenameErrors(siriusError model.ValidationErrors) model.ValidationErrors {
	mappedErrors := model.ValidationErrors{}
	for fieldName, value := range siriusError {
		for errorType, errorMessage := range value {
			err := make(map[string]string)
			if mapping, ok := validationMappings[fieldName][errorType]; ok {
				err[errorType] = mapping.v
				mappedErrors[mapping.k] = err
			} else {
				err[errorType] = errorMessage
				mappedErrors[fieldName] = err
			}
		}
	}
	return mappedErrors
}
