package main

import (
	"fmt"
	"time"
)

const (
	dateLayout = "2006-01-02"
)

type AwbDate struct {
	RenderDate
	Restmuell int `json:"restmuell,string"`
	Bio       int `json:"bio,string"`
	Wert      int `json:"wert,string"`
	Papier    int `json:"papier,string"`
}

type RenderDate struct {
	Termin    string `json:"termin"`
	Wochentag string `json:"wochentag"`
	Typ       string `json:"typ"`
}

func (date AwbDate) getDate() time.Time {
	time, err := time.Parse(dateLayout, date.Termin)

	if err != nil {
		fmt.Println(err)
	}

	return time
}

func (date AwbDate) isInFuture() bool {
	now := time.Now()
	return date.getDate().After(now)
}

func (date *AwbDate) setType() {
	date.Typ = "unbekannt"

	if date.Restmuell > 0 {
		date.Typ = "Restmüll"
	}
	if date.Bio > 0 {
		date.Typ = "Bio"
	}
	if date.Wert > 0 {
		date.Typ = "Wert"
	}
	if date.Papier > 0 {
		date.Typ = "Papier"
	}
}
