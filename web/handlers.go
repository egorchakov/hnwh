package web

import (
	"errors"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/evgorchakov/hnwh/Godeps/_workspace/src/github.com/Sirupsen/logrus"
	"github.com/evgorchakov/hnwh/database"
	"github.com/evgorchakov/hnwh/models"
)

type processedComment struct {
	Id   int
	By   string
	Time time.Time
	Text template.HTML
}

type searchBarData struct {
	Keywords string
	Months   int
}

type commentSearchData struct {
	SearchBarData searchBarData
	Comments      []processedComment
}

var (
	log = logrus.New()

	templates = template.Must(template.ParseGlob("templates/*html"))

	monthFieldError   = errors.New("Month value is invalid.")
	keywordFieldError = errors.New("Keywords value is invalid.")
	oopsError         = errors.New("Something went wrong!")
	queryError        = errors.New(
		"Provided keyword query is invalid. Here are examples of valid queries:\n" +
			"\n(python | golang) & London\n\n" +
			"haskell & 'san francisco'")

	StaticHandler = http.StripPrefix("/static/", http.FileServer(http.Dir("./static/")))
)

func processKeywords(keywords string) string {
	if strings.Count(keywords, `"`)%2 == 0 {
		return strings.Replace(keywords, `"`, `'`, -1)
	}

	return keywords
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "index", nil)
}

func CommentSearchHandler(w http.ResponseWriter, r *http.Request) {
	var (
		comments []models.HNComment
		err      error
	)

	query := r.URL.Query()

	if len(query) == 0 {
		IndexHandler(w, r)
		return
	}

	keywords := query.Get("keywords")
	if len(keywords) == 0 {
		http.Error(w, keywordFieldError.Error(), http.StatusBadRequest)
		log.Error(keywordFieldError)
		return
	}

	processedKeywords := processKeywords(keywords)

	months, err := strconv.Atoi(r.URL.Query().Get("months"))
	if err != nil {
		http.Error(w, monthFieldError.Error(), http.StatusBadRequest)
		log.Error(err)
		return
	}

	latestStories, err := database.GetLatestStories(months)
	if err != nil {
		http.Error(w, oopsError.Error(), http.StatusInternalServerError)
		log.Error(err)
		return
	}

	storyIDs := make([]int, len(latestStories))
	for i, story := range latestStories {
		storyIDs[i] = story.Id
	}

	comments, err = database.GetCommentsByKeywords(processedKeywords, storyIDs)
	if err != nil {
		http.Error(w, queryError.Error(), http.StatusInternalServerError)
		log.Error(err)
		return
	}

	processedComments := make([]processedComment, len(comments))
	for i, comment := range comments {
		processedComments[i] = processedComment{
			Id:   comment.Id,
			By:   comment.By,
			Time: comment.Time,
			Text: template.HTML(comment.Text),
		}
	}

	data := commentSearchData{
		SearchBarData: searchBarData{
			Keywords: keywords,
			Months:   months,
		},
		Comments: processedComments,
	}

	templates.ExecuteTemplate(w, "comments", data)
}
