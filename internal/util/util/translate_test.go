package util

import (
	"github.com/opg-sirius-finance-admin/internal/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRenameErrors(t *testing.T) {
	siriusErrors := model.ValidationErrors{
		"ReportType": map[string]string{"required": ""},
		"Date":       map[string]string{"Date": ""},
	}
	expected := model.ValidationErrors{
		"ReportType": map[string]string{"required": "Please select a report type"},
		"Date":       map[string]string{"Date": "Please select the report date"},
	}

	assert.Equal(t, expected, RenameErrors(siriusErrors))
}

func TestRenameErrors_default(t *testing.T) {
	siriusErrors := model.ValidationErrors{
		"x": map[string]string{"y": "z"},
	}

	assert.Equal(t, siriusErrors, RenameErrors(siriusErrors))
}
