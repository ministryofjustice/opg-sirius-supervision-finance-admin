package shared

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

type testJsonDateStruct struct {
	TestDate Date `json:"testDate"`
}

func TestDate_Before_And_After(t *testing.T) {
	tests := []struct {
		name              string
		date1             Date
		date2             Date
		wantForBeforeTest bool
		wantForAfterTest  bool
	}{
		{
			name:              "Date1 is before Date2",
			date1:             NewDate("01/01/2020"),
			date2:             NewDate("02/01/2020"),
			wantForBeforeTest: true,
			wantForAfterTest:  false,
		},
		{
			name:              "Date1 is after Date2",
			date1:             NewDate("02/01/2020"),
			date2:             NewDate("01/01/2020"),
			wantForBeforeTest: false,
			wantForAfterTest:  true,
		},
		{
			name:              "Date1 is the same as Date2",
			date1:             NewDate("01/01/2020"),
			date2:             NewDate("01/01/2020"),
			wantForBeforeTest: false,
			wantForAfterTest:  false,
		},
		{
			name:              "Date1 is empty",
			date1:             Date{},
			date2:             NewDate("02/01/2020"),
			wantForBeforeTest: true,
			wantForAfterTest:  false,
		},
		{
			name:              "Date2 is empty",
			date1:             NewDate("01/01/2020"),
			date2:             Date{},
			wantForBeforeTest: false,
			wantForAfterTest:  true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.wantForBeforeTest, test.date1.Before(test.date2))
			assert.Equal(t, test.wantForAfterTest, test.date1.After(test.date2))
		})
	}
}

func TestDate_IsNull(t *testing.T) {
	tests := []struct {
		name string
		date Date
		want bool
	}{
		{
			name: "Date passed in is not null",
			date: NewDate("01/01/2020"),
			want: false,
		},
		{
			name: "Date passed matches a nil date",
			date: NewDate("01/01/0001"),
			want: true,
		},
		{
			name: "Date passed is null",
			date: Date{},
			want: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.want, test.date.IsNull())
		})
	}
}

func TestDate_MarshalJSON(t *testing.T) {
	v := testJsonDateStruct{TestDate: NewDate("01/01/2020")}
	b, err := json.Marshal(v)
	assert.Nil(t, err)
	assert.Equal(t, `{"testDate":"01\/01\/2020"}`, string(b))
}

func TestDate_String(t *testing.T) {
	tests := []struct {
		name      string
		inputDate string
		want      string
	}{
		{
			name:      "returns correct format for slashers",
			inputDate: "01/01/2020",
			want:      "01/01/2020",
		},
		{
			name:      "returns correct format for dashers",
			inputDate: "2024-10-01",
			want:      "01/10/2024",
		},
		{
			name:      "returns correct format for date time string",
			inputDate: "2025-01-02T18:07:10+00:00",
			want:      "02/01/2025",
		},
		{
			name:      "returns nothing for default date string",
			inputDate: "01/01/0001",
			want:      "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewDate(tt.inputDate).String(), "stringToTime(%v)", tt.inputDate)
		})
	}
}

func TestDate_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		json string
		want string
	}{
		{
			json: `{"testDate":"01\/01\/2020"}`,
			want: "01/01/2020",
		},
		{
			json: `{"testDate":"01/01/2020"}`,
			want: "01/01/2020",
		},
		{
			json: `{"testDate":"2020-01-01T20:01:02+00:00"}`,
			want: "01/01/2020",
		},
	}
	for i, test := range tests {
		t.Run("Scenario "+strconv.Itoa(i+1), func(t *testing.T) {
			var v *testJsonDateStruct
			err := json.Unmarshal([]byte(test.json), &v)
			assert.Nil(t, err)
			assert.Equal(t, test.want, v.TestDate.String())
		})
	}
}
