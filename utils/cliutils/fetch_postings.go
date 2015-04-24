package main

import (
	"github.com/caser/gophernews"
	"github.com/evgorchakov/hnwh/models"
	"github.com/evgorchakov/hnwh/utils"
	"github.com/evgorchakov/hnwh/utils/dbutils"
	"golang.org/x/text/unicode/norm"
	"log"
	"os"
	"strconv"
	"time"
)

var (
	hn_post    *models.HNPost
	hn_story   *models.HNStory
	hn_comment *models.HNComment
)

func main() {
	// Read and convert the args
	story_args := os.Args[1:]
	story_ids := make([]int, len(story_args))

	for i := range story_args {
		k, err := strconv.Atoi(story_args[i])
		utils.FatalErr(err)
		story_ids[i] = k
	}

	// Init the client
	client := gophernews.NewClient()

	// Loop over stories and insert
	for _, story_id := range story_ids {

		story, err := client.GetStory(story_id)
		utils.FatalErr(err)

		log.Println("Processing story", story.ID)

		hn_post = &models.HNPost{
			Id:   story.ID,
			By:   story.By,
			Time: time.Unix(int64(story.Time), 0),
		}

		hn_story = &models.HNStory{
			HNPost: *hn_post,
			Score:  story.Score,
			Title:  story.Title,
			URL:    story.URL,
		}

		err = dbutils.InsertOrUpdateStory(hn_story)
		if err != nil {
			log.Println("Failed to insert story", story.ID, ":", err)
		}

		// Loop over the story's comments and insert
		for _, comment_id := range story.Kids {
			comment, err := client.GetComment(comment_id)
			log.Println("Processing comment", comment.ID)

			if err != nil {
				log.Println("Failed to get comment", comment_id, ":", err)
				continue
			}

			hn_post = &models.HNPost{
				Id:   comment.ID,
				By:   comment.By,
				Time: time.Unix(int64(comment.Time), 0),
			}

			hn_comment = &models.HNComment{
				HNPost:   *hn_post,
				Text:     norm.NFKC.String(comment.Text),
				ParentId: story.ID,
			}

			err = dbutils.InsertOrUpdateComment(hn_comment)
			if err != nil {
				log.Println("Failed to insert comment", story.ID, ":", err)
			}
		}
	}
}
