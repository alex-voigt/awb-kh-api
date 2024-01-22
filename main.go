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

	RenderDate struct {
		Date    string `json:"termin"`
		Weekday string `json:"wochentag"`
		Type    string `json:"typ"`
	}
)

func (date Calendar) getDate() time.Time {
	return time.UnixMilli(date.Date)
}

func (date HolidayView) getDate() time.Time {
	return time.UnixMilli(date.Holiday)
}

func (date HolidayView) getShiftDate() time.Time {
	return time.UnixMilli(date.ShiftTo)
}

func translateWeekday(weekday string) string {
	r := strings.NewReplacer(
		"Monday", "Montag",
		"Tuesday", "Dienstag",
		"Wednesday", "Mittwoch",
		"Thursday", "Donnerstag",
		"Friday", "Freitag",
		"Saturday", "Samstag",
		"Sunday", "Sonntag",
	)

	return r.Replace(weekday)
}

func (r AwbData) calculateNextDates(calendar Calendar, startDate time.Time) []RenderDate {
	until := startDate.AddDate(0, 0, DAYS_INTERVAL)
	renderDates := []RenderDate{}

	date := calendar.getDate()
	for date.Before(until) {
		shiftedDate := r.shiftHolidays(date)
		renderDate := RenderDate{
			Type:    calendar.Type,
			Date:    shiftedDate.Format(DATE_LAYOUT),
			Weekday: translateWeekday(shiftedDate.Weekday().String()),
		}
		renderDates = append(renderDates, renderDate)
		date = date.AddDate(0, 0, calendar.Frequency)
	}

	return renderDates
}

func (a AwbData) shiftHolidays(date time.Time) time.Time {
	for _, holidayView := range a.HolidayViews {
		if holidayView.IsActive && holidayView.getDate() == date {
			return holidayView.getShiftDate()
		}
	}

	return date
}

func getDatesSortedAsc(renderDates []RenderDate) []RenderDate {
	sort.Slice(renderDates, func(i, j int) bool {
		return renderDates[i].Date < renderDates[j].Date
	})

	return renderDates
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

func getDatesFromResponse(response AwbResponse, startDate time.Time) []RenderDate {
	renderDates := []RenderDate{}

	for _, calendar := range response.Data.Calendars {
		nextDates := response.Data.calculateNextDates(calendar, startDate)
		renderDates = append(renderDates, nextDates...)
	}

	sortedDates := getDatesSortedAsc(renderDates)

	if len(sortedDates) > MAX_ENTRIES {
		return sortedDates[0:MAX_ENTRIES]
	}
	return sortedDates
}

func getNextDates(w http.ResponseWriter, req *http.Request) {
	awbDates := postRequest(req)
	yesterday := time.Now().Add(-time.Hour * 24)
	responseDates := getDatesFromResponse(awbDates, yesterday)

	marshalled, err := json.Marshal(responseDates)
	if err != nil {
		log.Println(err)
	}

	w.Header().Add("Content-Type", "application/json")

	fmt.Fprintf(w, string(marshalled))
}

func main() {
	fmt.Println(fmt.Sprintf("starting web server on port %d", PORT))
	http.HandleFunc("/getDates", getNextDates)
	http.ListenAndServe(fmt.Sprintf(":%d", PORT), nil)
}
