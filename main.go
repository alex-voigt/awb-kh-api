package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/goodsign/monday"
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
		Calendars []AwbDate `json:"calendars"`
	}

	AwbDate struct {
		Type string `json:"name"`
		Date int64  `json:"fromDate"`
	}

	RenderDate struct {
		Date    string `json:"termin"`
		Weekday string `json:"wochentag"`
		Type    string `json:"typ"`
	}
)

func (a AwbResponse) getDatesSortedAsc() []AwbDate {
	dates := a.Data.Calendars

	sort.Slice(dates, func(i, j int) bool {
		return dates[i].Date < dates[j].Date
	})

	return dates
}

func (a AwbDate) toRenderDate() RenderDate {
	date := a.getDate()

	return RenderDate{
		Type:    a.Type,
		Date:    date.Format(DATE_LAYOUT),
		Weekday: monday.GetLongDays(monday.LocaleDeDE)[date.Weekday()],
	}
}

func (date AwbDate) getDate() time.Time {
	return time.UnixMilli(date.Date)
}

func (date AwbDate) isNextDate() bool {
	return date.isToday() || date.isInFuture()
}

func (date AwbDate) isToday() bool {
	y1, m1, d1 := time.Now().Date()
	y2, m2, d2 := date.getDate().Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

func (date AwbDate) isInFuture() bool {
	now := time.Now()
	return date.getDate().After(now)
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

func postRequest(req *http.Request) []AwbDate {
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

	dates := AwbResponse{}
	err = json.Unmarshal(body, &dates)
	if err != nil {
		log.Println(err)
	}

	return dates.getDatesSortedAsc()
}

func getProcessedDates(w http.ResponseWriter, req *http.Request) {
	awbDates := postRequest(req)

	entries := 0
	responseDates := []RenderDate{}
	for _, awbDate := range awbDates {
		if awbDate.isNextDate() {
			entries++
			responseDates = append(responseDates, awbDate.toRenderDate())
		}

		if entries == MAX_ENTRIES {
			break
		}
	}

	marshalled, err := json.Marshal(responseDates)
	if err != nil {
		log.Println(err)
	}

	w.Header().Add("Content-Type", "application/json")

	fmt.Fprintf(w, string(marshalled))
}

func main() {
	fmt.Println(fmt.Sprintf("starting web server on port %d", PORT))
	http.HandleFunc("/getDates", getProcessedDates)
	http.ListenAndServe(fmt.Sprintf(":%d", PORT), nil)
}
