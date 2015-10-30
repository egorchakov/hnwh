package hnutils

import (
	"errors"
	"fmt"
	"html"
	"sort"
	"time"

	"github.com/evgorchakov/hnwh/Godeps/_workspace/src/github.com/Sirupsen/logrus"
	hn "github.com/evgorchakov/hnwh/Godeps/_workspace/src/github.com/caser/gophernews"

	"github.com/evgorchakov/hnwh/database"
	"github.com/evgorchakov/hnwh/models"
)

const (
	author = "whoishiring"
	title  = "Ask HN: Who is hiring? (%s)"
)

var (
	log    = logrus.New()
	client hn.Client
)

func getDesiredStoryTitle(t time.Time) string {
	date := fmt.Sprintf("%s %d", t.Month(), t.Year())
	return fmt.Sprintf(title, date)
}

func getHiringStory(submissionTime time.Time) (hn.Story, error) {
	var (
		story hn.Story
		err   error
	)

	user, err := client.GetUser(author)
	if err != nil {
		return story, err
	}

	userPostIds := user.Submitted
	sort.Sort(sort.Reverse(sort.IntSlice(userPostIds)))
	storyTitle := getDesiredStoryTitle(submissionTime)

	for _, post_id := range userPostIds {
		story, _ = client.GetStory(post_id)
		if story.Title == storyTitle {
			return story, nil
		}
	}
	return story, errors.New("No \"Who's Hiring?\" story found")
}

func getTopLevelCommentsForStory(story *hn.Story) []hn.Comment {
	var (
		comments []hn.Comment
		comment  hn.Comment
		err      error
	)

	for _, comment_id := range story.Kids {
		comment, err = client.GetComment(comment_id)
		if err == nil {
			comments = append(comments, comment)
		} else {
			log.Error(err)
		}
	}
	return comments
}

func convertStoryToModel(story *hn.Story) models.HNStory {
	storyTime := time.Unix(int64(story.Time), 0)
	hnPost := models.HNPost{
		Id:   story.ID,
		By:   story.By,
		Time: storyTime,
	}
	hnStory := models.HNStory{
		HNPost: hnPost,
		Score:  story.Score,
		Title:  story.Title,
		URL:    story.URL,
	}
	return hnStory
}

func convertCommentToModel(comment *hn.Comment, storyId int) models.HNComment {
	comment_time := time.Unix(int64(comment.Time), 0)
	hnPost := models.HNPost{
		Id:   comment.ID,
		By:   comment.By,
		Time: comment_time,
	}
	hnComment := models.HNComment{
		HNPost:   hnPost,
		Text:     html.UnescapeString(comment.Text),
		ParentId: storyId,
	}
	return hnComment
}

func UpdateComments(HNClient *hn.Client, date time.Time) {
	var err error
	client = *HNClient

	story, err := getHiringStory(date)
	if err != nil {
		log.Fatal(err)
	}

	comments := getTopLevelCommentsForStory(&story)
	if len(comments) == 0 {
		log.Fatal("No comments found")
	}

	storyModel := convertStoryToModel(&story)
	err = database.InsertOrUpdateStory(&storyModel)
	if err != nil {
		log.Fatal(err)
	}

	var commentModel models.HNComment
	for _, comment := range comments {
		commentModel = convertCommentToModel(&comment, story.ID)
		err = database.InsertOrUpdateComment(&commentModel)
		if err != nil {
			log.Error(err)
		}
	}
}
