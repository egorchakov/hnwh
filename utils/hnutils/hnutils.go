package hnutils

import (
	"fmt"
	"log"
	"sort"
	"time"

	hn "github.com/caser/gophernews"
)

const (
	author = "whoishiring"
	title  = "Ask HN: Who is hiring? (%s)"
)

func getStoryTitle(t time.Time) string {
	date := fmt.Sprintf("%s %d", t.Month(), t.Year())
	return fmt.Sprintf(title, date)
}

func GetHiringStory(client *hn.Client, t time.Time) []hn.Story {
	user, err := client.GetUser(author)
	if err != nil {
		log.Fatalln(err.Error())
	}

	userStories := make([]hn.Story, 0)
	userStoryIds := user.Submitted
	sort.Sort(sort.Reverse(sort.IntSlice(userStoryIds)))

	storyTitle := getStoryTitle(t)
	for _, object_id := range userStoryIds {
		story, err := client.GetStory(object_id)
		if story.Title == storyTitle && err == nil {
			userStories = append(userStories, story)
		}
	}
	return userStories
}
