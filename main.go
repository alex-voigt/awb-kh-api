package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

const (
	PORT        = 8010
	API_URL     = "https://app.awb-bad-kreuznach.de/api/loadDates.php"
	MAX_ENTRIES = 10
)

func postAwbRequest(req *http.Request) []AwbDate {
	client := &http.Client{}
	request, err := http.NewRequest("POST", API_URL, strings.NewReader(req.URL.Query().Encode()))
	if err != nil {
		log.Fatal(err)
	}
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	response, err := client.Do(request)
	if err != nil {
		log.Println(err)
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
	}

	dates := []AwbDate{}
	json.Unmarshal(body, &dates)

	return dates
}

func getDates(w http.ResponseWriter, req *http.Request) {
	awbDates := postAwbRequest(req)

	entries := 0
	responseDates := []RenderDate{}
	for _, awbDate := range awbDates {
		awbDate.setType()
		if awbDate.isInFuture() {
			entries++
			responseDates = append(responseDates, awbDate.RenderDate)
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
	http.HandleFunc("/getDates", getDates)
	http.ListenAndServe(fmt.Sprintf(":%d", PORT), nil)
}
