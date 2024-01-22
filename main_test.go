package main

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"
	"time"
)

func TestDateFiltering(t *testing.T) {
	awbResponseFile, _ := os.ReadFile("fixtures/response_january.json")
	var response AwbResponse
	json.Unmarshal(awbResponseFile, &response)

	startTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.Local)
	dates := getDatesFromResponse(response, startTime)

	expectedDates := []RenderDate{
		RenderDate{
			Date:    "2024-01-04",
			Weekday: "Donnerstag",
			Type:    "Restabfall",
		},
		RenderDate{
			Date:    "2024-01-05",
			Weekday: "Freitag",
			Type:    "Gelber Sack/Gelbe Tonne",
		},
		RenderDate{
			Date:    "2024-01-10",
			Weekday: "Mittwoch",
			Type:    "Bioabfall",
		},
		RenderDate{
			Date:    "2024-01-16",
			Weekday: "Dienstag",
			Type:    "Altpapier",
		},
		RenderDate{
			Date:    "2024-01-17",
			Weekday: "Mittwoch",
			Type:    "Restabfall",
		},
		RenderDate{
			Date:    "2024-01-24",
			Weekday: "Mittwoch",
			Type:    "Bioabfall",
		},
	}

	if len(dates) != len(expectedDates) {
		t.Errorf("got %d dates; expected %d", len(dates), len(expectedDates))
	}

	for i, expectedDate := range expectedDates {
		if !reflect.DeepEqual(expectedDate, dates[i]) {
			t.Errorf("got %+v, expected: %+v, received und expected date does not match!", dates[i], expectedDate)
		}
	}
}
