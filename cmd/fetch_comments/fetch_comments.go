package main

import (
	"log"
	"time"

	"github.com/evgorchakov/hnwh/Godeps/_workspace/src/github.com/caser/gophernews"
	"github.com/evgorchakov/hnwh/Godeps/_workspace/src/gopkg.in/alecthomas/kingpin.v2"
	"github.com/evgorchakov/hnwh/database"
	"github.com/evgorchakov/hnwh/utils/hnutils"
)

const (
	dateFormat = "January 2006"
)

var (
	date = kingpin.Flag("date", "Month and year of the posting, e.g. October 2015").String()
)

func main() {
	kingpin.Parse()
	date, err := time.Parse(dateFormat, *date)
	if err != nil {
		log.Fatal(err)
	}

	database.SetupDB()
	database.InitDB()
	client := gophernews.NewClient()
	hnutils.UpdateComments(client, date)
}
