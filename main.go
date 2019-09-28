package main

import (
	"flag"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gen2brain/beeep"
	"github.com/gocolly/colly"
	"github.com/lestrrat-go/strftime"
)

type earthquake struct {
	Date      string
	Time      string
	DateTime  time.Time
	Latitude  float64
	Longitude float64
	Depth     float64
	Md        float64
	Ml        float64
	Mw        float64
	Yer       string
}

var (
	earthquakes []earthquake
	day         string
	hour        string
	minute      string
	duration    time.Duration
	min         float64
)

func init() {
	day, _ = strftime.Format(`%Y.%m.%d`, time.Now())
	hour, _ = strftime.Format(`%H`, time.Now())
	minute, _ = strftime.Format(`%M`, time.Now())

	var period int
	flag.IntVar(&period, "period", 5, "periodic interval")

	var min int
	flag.IntVar(&min, "min", 3, "min")
	flag.Parse()

	duration = time.Minute * time.Duration(period)
}

func main() {
	log.Printf("[*] BDTİM sayfası her %s de bir kontrol edilecek.", duration)
	run()
	select {}
}

func run() {
	c := colly.NewCollector()
	c.OnHTML("body > pre", func(e *colly.HTMLElement) {
		table := strings.Split(e.Text, "--------------")
		rows := strings.Split(table[len(table)-1], "\n")
		for _, row := range rows {
			if len(row) > 0 {
				parsed := regexp.MustCompile(`\s+`).Split(row, 10)

				latitude, _ := strconv.ParseFloat(parsed[2], 64)
				longitude, _ := strconv.ParseFloat(parsed[3], 64)
				depth, _ := strconv.ParseFloat(parsed[4], 64)
				md, _ := strconv.ParseFloat(parsed[5], 64)
				ml, _ := strconv.ParseFloat(parsed[6], 64)
				mw, _ := strconv.ParseFloat(parsed[7], 64)

				dots := regexp.MustCompile(`\.`)
				dateString := dots.ReplaceAllString(parsed[0], "-") + "T" + parsed[1] + "+30:00"
				dateTime, _ := time.Parse(time.RFC3339, dateString)

				spaces := regexp.MustCompile(`\s+`)
				d := earthquake{
					Date:      parsed[0],
					Time:      parsed[1],
					DateTime:  dateTime,
					Latitude:  latitude,
					Longitude: longitude,
					Depth:     depth,
					Md:        md,
					Ml:        ml,
					Mw:        mw,
					Yer:       spaces.ReplaceAllString(parsed[8], " "),
				}
				earthquakes = append(earthquakes, d)
			}
		}
	})

	log.Print("Sayfa ziyaret ediliyor.")
	c.Visit("http://www.koeri.boun.edu.tr/scripts/lst4.asp")

	last := last(earthquakes)
	biggers := biggers(last)

	for _, x := range biggers {
		beeep.Notify(
			"BDTİM Uyarı "+x.Date+" "+x.Time,
			fmt.Sprintf("%s de %v şiddetinde earthquake oldu.", x.Yer, x.Ml),
			"assets/logo.gif",
		)
	}

	time.AfterFunc(duration, run)
}

func last(list []earthquake) []earthquake {
	newlist := list[:0]
	now := time.Now()
	for _, x := range list {
		if duration.Minutes() >= now.Sub(x.DateTime).Minutes() {
			newlist = append(newlist, x)
		}
	}

	return newlist
}

func biggers(list []earthquake) []earthquake {
	newlist := list[:0]
	for _, x := range list {
		if x.Ml > min {
			newlist = append(newlist, x)
		}
	}

	return newlist
}
