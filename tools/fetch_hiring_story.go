package main

import (
	"html"
	"log"
	"time"

	"github.com/caser/gophernews"
	"github.com/evgorchakov/hnwh/database"
	"github.com/evgorchakov/hnwh/models"
	"github.com/evgorchakov/hnwh/utils/hnutils"
	"gopkg.in/alecthomas/kingpin.v2"
)

const (
	dateFormat = "January 2006"
)

var (
	hn_post    *models.HNPost
	hn_story   *models.HNStory
	hn_comment *models.HNComment

	date = kingpin.Flag("date", "Month and year of the posting, e.g. October 2015").String()
)

func main() {
	kingpin.Parse()
	date, err := time.Parse(dateFormat, *date)
	if err != nil {
		log.Fatal(err)
	}

	client := gophernews.NewClient()
	stories := hnutils.GetHiringStory(client, date)

	if len(stories) == 0 {
		log.Fatal("No stories found!")
	}

	database.InitDB()
	// Loop over stories and insert
	for _, story := range stories {
		story_time := time.Unix(int64(story.Time), 0)

		log.Printf("Processing story %+v by %+v posted on %+v", story.ID, story.By, story_time)

		hn_post = &models.HNPost{
			Id:   story.ID,
			By:   story.By,
			Time: story_time,
		}

		hn_story = &models.HNStory{
			HNPost: *hn_post,
			Score:  story.Score,
			Title:  story.Title,
			URL:    story.URL,
		}

		err := database.InsertOrUpdateStory(hn_story)
		if err != nil {
			log.Printf("Failed to insert story %+v: %+v", story.ID, err)
		}

		// Loop over the story's comments and insert
		for _, comment_id := range story.Kids {
			comment, err := client.GetComment(comment_id)
			comment_time := time.Unix(int64(comment.Time), 0)

			log.Printf("Processing comment %+v by %+v posted on %+v", comment.ID, comment.By, comment_time)

			if err != nil {
				log.Printf("Failed to get comment %+v: %+v", comment.ID, err)
				continue
			}

			hn_post = &models.HNPost{
				Id:   comment.ID,
				By:   comment.By,
				Time: comment_time,
			}

			hn_comment = &models.HNComment{
				HNPost:   *hn_post,
				Text:     html.UnescapeString(comment.Text),
				ParentId: story.ID,
			}

			err = database.InsertOrUpdateComment(hn_comment)
			if err != nil {
				log.Printf("Failed to insert comment %+v: %+v", comment.ID, err)
			}
		}
	}
}
