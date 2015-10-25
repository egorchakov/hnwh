package web

import (
	"errors"
	"net/http"
	"strconv"
	"text/template"

	"github.com/Sirupsen/logrus"

	"github.com/evgorchakov/hnwh/database"
	"github.com/evgorchakov/hnwh/models"
)

type searchBarData struct {
	Keywords string
	Months   int
}
type commentSearchData struct {
	SearchBarData searchBarData
	Comments      []models.HNComment
}

const (
	keywordsPlaceholder = "(python | 'embedded systems') & london"
	monthsPlaceholder   = 2
)

var (
	log = logrus.New()

	templates = template.Must(template.ParseGlob("templates/*html"))

	monthFieldError   = errors.New("Month value is invalid.")
	keywordFieldError = errors.New("Keywords value is invalid.")
	oopsError         = errors.New("Something went wrong!")

	StaticHandler = http.StripPrefix("/static/", http.FileServer(http.Dir("./static/")))
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	data := commentSearchData{
		SearchBarData: searchBarData{
			Keywords: keywordsPlaceholder,
			Months:   monthsPlaceholder,
		},
		Comments: nil,
	}

	templates.ExecuteTemplate(w, "comments", data)
}

func CommentSearchHandler(w http.ResponseWriter, r *http.Request) {
	templates = template.Must(template.ParseGlob("templates/*html"))
	var (
		comments []models.HNComment
		err      error
	)

	keywords := r.URL.Query().Get("keywords")
	if len(keywords) == 0 {
		http.Error(w, keywordFieldError.Error(), http.StatusBadRequest)
		log.Error(keywordFieldError)
		return
	}

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

	comments, err = database.GetCommentsByKeywords(keywords, storyIDs)
	if err != nil {
		http.Error(w, oopsError.Error(), http.StatusInternalServerError)
		log.Error(err)
		return
	}

	data := commentSearchData{
		SearchBarData: searchBarData{
			Keywords: keywords,
			Months:   months,
		},
		Comments: comments,
	}

	templates.ExecuteTemplate(w, "comments", data)
}
