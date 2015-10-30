package main

import (
	"net/http"
	"os"

	"github.com/evgorchakov/hnwh/Godeps/_workspace/src/github.com/Sirupsen/logrus"
	"github.com/evgorchakov/hnwh/Godeps/_workspace/src/github.com/gorilla/mux"
	"github.com/evgorchakov/hnwh/database"
	"github.com/evgorchakov/hnwh/web"
)

var log = logrus.New()

func main() {
	log.Info("Here we go!")
	database.SetupDB()

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", web.IndexHandler)
	router.HandleFunc("/search", web.CommentSearchHandler)
	router.PathPrefix("/static/").Handler(web.StaticHandler)
	http.Handle("/", router)
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), router))
}
