package main

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"
	"time"
)

func TestWasteDateCalculation(t *testing.T) {
	awbResponseFile, _ := os.ReadFile("fixtures/response_january.json")
	var response AwbResponse
	json.Unmarshal(awbResponseFile, &response)

	startTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.Local)
	dates := response.getWasteDates(startTime)

	expectedDates := []WasteDate{
		WasteDate{
			Date: time.Date(2024, 1, 4, 0, 0, 0, 0, time.Local),
			Type: "Restabfall",
		},
		WasteDate{
			Date: time.Date(2024, 1, 5, 0, 0, 0, 0, time.Local),
			Type: "Gelber Sack/Gelbe Tonne",
		},
		WasteDate{
			Date: time.Date(2024, 1, 10, 0, 0, 0, 0, time.Local),
			Type: "Bioabfall",
		},
		WasteDate{
			Date: time.Date(2024, 1, 16, 0, 0, 0, 0, time.Local),
			Type: "Altpapier",
		},
		WasteDate{
			Date: time.Date(2024, 1, 17, 0, 0, 0, 0, time.Local),
			Type: "Restabfall",
		},
		WasteDate{
			Date: time.Date(2024, 1, 24, 0, 0, 0, 0, time.Local),
			Type: "Bioabfall",
		},
	}

	if len(dates) == 0 {
		t.Errorf("no waste dates found")
	}

	if len(dates) != len(expectedDates) {
		t.Errorf("got %d dates; expected %d", len(dates), len(expectedDates))
	}

	for i, date := range dates {
		if !reflect.DeepEqual(date, expectedDates[i]) {
			t.Errorf("index %d: got %+v, expected: %+v, received und expected date does not match!", i, date, expectedDates[i])
		}
	}
}

func TestWasteDateFiltering(t *testing.T) {
	awbResponseFile, _ := os.ReadFile("fixtures/response_january.json")
	var response AwbResponse
	json.Unmarshal(awbResponseFile, &response)

	startTime := time.Date(2024, 1, 17, 0, 0, 0, 0, time.Local)
	dates := response.getWasteDates(startTime)

	expectedDates := []WasteDate{
		WasteDate{
			Type: "Restabfall",
			Date: time.Date(2024, 1, 17, 0, 0, 0, 0, time.Local),
		},
		WasteDate{
			Type: "Bioabfall",
			Date: time.Date(2024, 1, 24, 0, 0, 0, 0, time.Local),
		},
		WasteDate{
			Type: "Restabfall",
			Date: time.Date(2024, 1, 31, 0, 0, 0, 0, time.Local),
		},
		WasteDate{
			Type: "Gelber Sack/Gelbe Tonne",
			Date: time.Date(2024, 2, 1, 0, 0, 0, 0, time.Local),
		},
		WasteDate{
			Type: "Bioabfall",
			Date: time.Date(2024, 2, 7, 0, 0, 0, 0, time.Local),
		},
		WasteDate{
			Type: "Altpapier",
			Date: time.Date(2024, 2, 13, 0, 0, 0, 0, time.Local),
		},
		WasteDate{
			Type: "Restabfall",
			Date: time.Date(2024, 2, 14, 0, 0, 0, 0, time.Local),
		},
	}

	if len(dates) == 0 {
		t.Errorf("no waste dates found")
	}

	if len(dates) != len(expectedDates) {
		t.Errorf("got %d dates; expected %d", len(dates), len(expectedDates))
	}

	for i, date := range dates {
		if !reflect.DeepEqual(date, expectedDates[i]) {
			t.Errorf("index %d: got %+v, expected: %+v, received und expected date does not match!", i, date, expectedDates[i])
		}
	}
}

func TestWasteDateToRenderDate(t *testing.T) {
	wasteDates := []WasteDate{
		WasteDate{
			Type: "Restabfall",
			Date: time.Date(2024, 1, 17, 0, 0, 0, 0, time.Local),
		},
		WasteDate{
			Type: "Bioabfall",
			Date: time.Date(2024, 1, 24, 0, 0, 0, 0, time.Local),
		},
		WasteDate{
			Type: "Restabfall",
			Date: time.Date(2024, 1, 31, 0, 0, 0, 0, time.Local),
		},
		WasteDate{
			Type: "Gelber Sack/Gelbe Tonne",
			Date: time.Date(2024, 2, 1, 0, 0, 0, 0, time.Local),
		},
	}

	renderDates := getRenderDates(wasteDates, 2)

	expectedRenderDates := []RenderDate{
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

	if len(renderDates) != len(expectedRenderDates) {
		t.Errorf("wrong amount of waste dates returned. got: %d, expected: %d", len(renderDates), len(expectedRenderDates))
	}

	if !reflect.DeepEqual(renderDates, expectedRenderDates) {
		t.Errorf("got %+v, expected: %+v, received und expected date does not match!", renderDates, expectedRenderDates)
	}
}

func TestGetRenderDatesWithLimit(t *testing.T) {
	wasteDate := WasteDate{
		Date: time.Date(2024, 1, 24, 0, 0, 0, 0, time.Local),
		Type: "Bioabfall",
	}

	renderDate := wasteDate.toRenderDate()

	expectedRenderDate := RenderDate{
		Date:    "2024-01-24",
		Weekday: "Mittwoch",
		Type:    "Bioabfall",
	}

	if !reflect.DeepEqual(renderDate, expectedRenderDate) {
		t.Errorf("got %+v, expected: %+v, received und expected date does not match!", renderDate, expectedRenderDate)
	}
}
