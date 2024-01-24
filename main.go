package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	PORT          = 8010
	API_URL       = "https://blupassionsystem.de/city/rest/garbageorte/getAllGarbageCalendar"
	MAX_ENTRIES   = 10
	DAYS_INTERVAL = 30
	DATE_LAYOUT   = "2006-01-02"
)

type (
	AwbResponse struct {
		Data AwbData `json:"data"`
	}

	AwbData struct {
		Calendars    []Calendar    `json:"calendars"`
		HolidayViews []HolidayView `json:"holidayViews"`
	}

	Calendar struct {
		Type      string `json:"name"`
		Date      int64  `json:"fromDate"`
		Frequency int    `json:"frequency"`
	}

	HolidayView struct {
		IsActive bool  `json:"active"`
		Holiday  int64 `json:"holiday"`
		ShiftTo  int64 `json:"shiftTo"`
	}

	WasteDate struct {
		Type string
		Date time.Time
	}

	RenderDate struct {
		Date    string `json:"termin"`
		Weekday string `json:"wochentag"`
		Type    string `json:"typ"`
	}
)

func (r AwbResponse) getWasteDates(startDate time.Time) []WasteDate {
	wasteDates := []WasteDate{}

	for _, calendar := range r.Data.Calendars {
		nextDates := r.Data.calculateDates(calendar, startDate)
		wasteDates = append(wasteDates, nextDates...)
	}

	return getDatesSortedAsc(wasteDates)
}

func getDatesSortedAsc(wasteDates []WasteDate) []WasteDate {
	sort.Slice(wasteDates, func(i, j int) bool {
		return wasteDates[i].Date.Before(wasteDates[j].Date)
	})

	return wasteDates
}

func (r AwbData) calculateDates(calendar Calendar, startDate time.Time) []WasteDate {
	until := startDate.AddDate(0, 0, DAYS_INTERVAL)
	wasteDates := []WasteDate{}

	date := calendar.getDate()
	for date.Before(until) {
		wasteDate := WasteDate{
			Type: calendar.Type,
			Date: r.getShiftedDate(date),
		}

		if wasteDate.isNextDate(startDate) {
			wasteDates = append(wasteDates, wasteDate)
		}
		date = date.AddDate(0, 0, calendar.Frequency)
	}

	return wasteDates
}

func (a AwbData) getShiftedDate(date time.Time) time.Time {
	for _, holidayView := range a.HolidayViews {
		if holidayView.IsActive && holidayView.getDate() == date {
			return holidayView.getShiftDate()
		}
	}

	return date
}

func (date Calendar) getDate() time.Time {
	return time.UnixMilli(date.Date)
}

func (date HolidayView) getDate() time.Time {
	return time.UnixMilli(date.Holiday)
}

func (date HolidayView) getShiftDate() time.Time {
	return time.UnixMilli(date.ShiftTo)
}

func (w WasteDate) isNextDate(startDate time.Time) bool {
	return w.isToday(startDate) || w.isInFuture(startDate)
}

func (w WasteDate) isToday(startDate time.Time) bool {
	y1, m1, d1 := startDate.Date()
	y2, m2, d2 := w.Date.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

func (w WasteDate) isInFuture(startDate time.Time) bool {
	return w.Date.After(startDate)
}

func (w WasteDate) toRenderDate() RenderDate {
	return RenderDate{
		Type:    w.Type,
		Date:    w.Date.Format(DATE_LAYOUT),
		Weekday: w.getTranslatedWeekDay(),
	}
}

func (w WasteDate) getTranslatedWeekDay() string {
	r := strings.NewReplacer(
		"Monday", "Montag",
		"Tuesday", "Dienstag",
		"Wednesday", "Mittwoch",
		"Thursday", "Donnerstag",
		"Friday", "Freitag",
		"Saturday", "Samstag",
		"Sunday", "Sonntag",
	)

	return r.Replace(w.Date.Weekday().String())
}

func buildRequest(req *http.Request) *http.Request {
	yesterday := time.Now().Add(-time.Hour * 24)
	fromTime := strconv.FormatInt(yesterday.UnixMilli(), 10)
	toTime := strconv.FormatInt(yesterday.AddDate(0, 0, DAYS_INTERVAL).UnixMilli(), 10)

	q := req.URL.Query()
	q.Add("fromTime", fromTime)
	q.Add("toTime", toTime)

	request, err := http.NewRequest("GET", fmt.Sprintf("%s?%s", API_URL, q.Encode()), nil)
	if err != nil {
		log.Fatal(err)
	}
	request.Header.Add("Accept", "application/json")
	request.Header.Add("Content-Type", "application/json")

	return request
}

func postRequest(req *http.Request) AwbResponse {
	client := &http.Client{}
	response, err := client.Do(buildRequest(req))
	if err != nil {
		log.Println(err)
	}

	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
	}

	awbResponse := AwbResponse{}
	err = json.Unmarshal(body, &awbResponse)
	if err != nil {
		log.Println(err)
	}

	return awbResponse
}

func getNextDates(w http.ResponseWriter, req *http.Request) {
	awbDates := postRequest(req)
	yesterday := time.Now().Add(-time.Hour * 24)
	wasteDates := awbDates.getWasteDates(yesterday)
	renderDates := getRenderDates(wasteDates, MAX_ENTRIES)

	marshalled, err := json.Marshal(renderDates)
	if err != nil {
		log.Println(err)
	}

	w.Header().Add("Content-Type", "application/json")

	fmt.Fprintf(w, string(marshalled))
}

func getRenderDates(wasteDates []WasteDate, maxEntries int) []RenderDate {
	renderDates := []RenderDate{}
	for _, wasteDate := range wasteDates {
		renderDates = append(renderDates, wasteDate.toRenderDate())
	}

	if len(renderDates) > maxEntries {
		renderDates = renderDates[0:maxEntries]
	}

	return renderDates
}

func main() {
	fmt.Println(fmt.Sprintf("starting web server on port %d", PORT))
	http.HandleFunc("/getDates", getNextDates)
	http.ListenAndServe(fmt.Sprintf(":%d", PORT), nil)
}
