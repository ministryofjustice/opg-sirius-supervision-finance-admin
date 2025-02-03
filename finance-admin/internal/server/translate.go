package server

import "github.com/ministryofjustice/opg-sirius-supervision-finance-admin/finance-admin/internal/model"

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
		"required":         pair{"Date", "Please select the report date"},
		"date-in-the-past": pair{"Date", "The report date must be today or in the past"},
	},
	"FromDate": {
		"FromDate": pair{"FromDate", "Date From must be before Date To"},
	},
	"ToDate": {
		"ToDate": pair{"ToDate", "Date To must be after Date From"},
	},
	"ReportUploadType": {
		"required": pair{"ReportUploadType", "Please select a report type"},
	},
	"UploadDate": {
		"required":         pair{"UploadDate", "Please enter a date"},
		"date-in-the-past": pair{"UploadDate", "The report date must be today or in the past"},
	},
	"Email": {
		"required": pair{"Email", "Please enter an email address to be notified when the upload is complete"},
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
