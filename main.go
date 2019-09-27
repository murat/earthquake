package main

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

type deprem struct {
	Tarih    string
	Saat     string
	Enlem    float64
	Boylam   float64
	Derinlik float64
	SiddetML float64
	SiddetMW float64
	SiddetMD float64
	Yer      string
}

type depremler struct {
	list []deprem
}

func main() {
	depremList := depremler{}
	c := colly.NewCollector()
	c.OnHTML("body > pre", func(e *colly.HTMLElement) {
		table := strings.Split(e.Text, "--------------")
		rows := strings.Split(table[len(table)-1], "\n")
		for _, row := range rows {
			if len(row) > 0 {
				parsed := regexp.MustCompile(`\s+`).Split(row, 10)
				enlem, _ := strconv.ParseFloat(strings.Join(parsed[3:4], ""), 64)
				boylam, _ := strconv.ParseFloat(strings.Join(parsed[4:5], ""), 64)
				derinlik, _ := strconv.ParseFloat(strings.Join(parsed[5:6], ""), 64)
				siddetmd, _ := strconv.ParseFloat(strings.Join(parsed[6:7], ""), 64)
				siddetml, _ := strconv.ParseFloat(strings.Join(parsed[7:8], ""), 64)
				siddetmw, _ := strconv.ParseFloat(strings.Join(parsed[8:9], ""), 64)

				spaces := regexp.MustCompile(`\s+`)
				d := deprem{
					Tarih:    fmt.Sprintf("%s", parsed[0]),
					Saat:     fmt.Sprintf("%s", parsed[1]),
					Enlem:    enlem,
					Boylam:   boylam,
					Derinlik: derinlik,
					SiddetMD: siddetmd,
					SiddetML: siddetml,
					SiddetMW: siddetmw,
					Yer:      spaces.ReplaceAllString(fmt.Sprintf("%s", parsed[9]), " "),
				}
				depremList.list = append(depremList.list, d)
			}
		}
	})

	c.Visit("http://www.koeri.boun.edu.tr/scripts/lst4.asp")

	fmt.Printf("%v", depremList.enBuyuk())
}

func (depremler depremler) enBuyuk() deprem {
	sort.Slice(depremler.list, func(i, j int) bool {
		return depremler.list[i].SiddetMD > depremler.list[j].SiddetMD
	})

	return depremler.list[0]
}
