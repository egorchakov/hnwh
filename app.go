package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/evgorchakov/hnwh/models"
	u "github.com/evgorchakov/hnwh/utils"
	"github.com/evgorchakov/hnwh/utils/dbutils"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

var (
	db   *sqlx.DB
	logs chan string
)

const (
	dateFormat = "2006-01-02"
)

func processURLQuery(q url.Values) (map[string]interface{}, error) {
	var params = map[string]interface{}{}

	tags := u.StrFilter(q["tag"], u.StrNotEmpty)

	if len(tags) == 0 {
		return params, errors.New("No tags provided")
	}

	logs <- fmt.Sprintf("Raw tags: %+v", tags)

	for i := range tags {
		tags[i] = strings.Replace(tags[i], " ", "+", -1)
	}

	logs <- fmt.Sprintf("Prepared tags: %+v", tags)

	params["tags"] = strings.Join(tags, "&")

	for _, param := range []string{"from", "to"} {
		val := q.Get(param)
		if val != "" {
			_, err := time.Parse(dateFormat, val)
			if err != nil {
				logs <- err.Error()
				return params, err
			}

			params[param] = val
		}
	}

	return params, nil
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	var comments []models.HNComment

	params, err := processURLQuery(r.URL.Query())
	if err != nil {
		logs <- err.Error()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logs <- fmt.Sprintf("Processed params: %+v", params)

	query := dbutils.BuildFullTextSearhQuery(params)
	logs <- query

	named_stmt, err := db.PrepareNamed(query)

	if err != nil {
		logs <- err.Error()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = named_stmt.Select(&comments, params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	js, err := json.Marshal(comments)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func commentHandler(w http.ResponseWriter, r *http.Request) {
	var comment models.HNComment
	vars := mux.Vars(r)

	query := "SELECT * FROM hn_comments WHERE id=$1"
	err := db.Get(&comment, query, vars["commentID"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	js, err := json.Marshal(comment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func logging(logs chan string) {
	for msg := range logs {
		log.Println(msg)
	}
}

func main() {
	logs = make(chan string, 1e5)
	go logging(logs)
	logs <- "Here we go!"

	dbutils.InitDB()
	db = dbutils.OpenDB()
	defer db.Close()
	logs <- "Initialized DB"

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/search", searchHandler)
	router.HandleFunc("/comments/{commentID}", commentHandler)

	log.Fatal(http.ListenAndServe("localhost:8000", router))
}
