package database

import (
	"strings"

	"github.com/evgorchakov/hnwh/Godeps/_workspace/src/github.com/jmoiron/sqlx"
	"github.com/evgorchakov/hnwh/models"
)

const (
	stmtGetStories                 = `SELECT * FROM hn_stories`
	stmtGetLatestStories           = `SELECT * from hn_stories ORDER BY time DESC LIMIT $1`
	stmtGetStoryById               = `SELECT * FROM hn_stories WHERE id = $1`
	stmtGetComments                = `SELECT * FROM hn_comments`
	stmtGetCommentById             = `SELECT * FROM hn_comments WHERE id = $1`
	stmtGetCommentsByKeywords      = `SELECT * FROM hn_comments WHERE parent IN (:parent_ids) AND text @@ to_tsquery(:keywords) ORDER BY time DESC`
	stmtGetCommentsByKeywordsPlain = `SELECT * FROM hn_comments WHERE parent IN (:parent_ids) AND text @@ plainto_tsquery(:keywords) ORDER BY time DESC`
)

func isPlainTSQuery(query string) bool {
	return !strings.ContainsAny(query, "'()&|!")
}

func GetLatestStories(count int) ([]models.HNStory, error) {
	var stories []models.HNStory
	err := db.Select(&stories, stmtGetLatestStories, count)
	if err != nil {
		log.Error("GetLatestStories failed: ", err)
		return nil, err
	}

	return stories, nil
}

func GetStories() ([]models.HNStory, error) {
	var stories []models.HNStory
	err := db.Select(&stories, stmtGetStories)

	if err != nil {
		log.Error("GetStories failed", err)
		return nil, err
	}

	return stories, nil
}

func GetStoryById(id int) (*models.HNStory, error) {
	var story models.HNStory
	err := db.Get(&story, stmtGetStoryById, id)

	if err != nil {
		log.Error("GetStoryById failed", err)
		return nil, err
	}

	return &story, nil
}

func GetComments() ([]models.HNComment, error) {
	var comments []models.HNComment
	err := db.Select(&comments, stmtGetComments)

	if err != nil {
		return nil, err
	}

	return comments, nil
}

func GetCommentById(id int) (*models.HNComment, error) {
	var comment models.HNComment
	err := db.Get(&comment, stmtGetCommentById, id)

	if err != nil {
		log.Error("GetCommentById failed", err)
		return nil, err
	}

	return &comment, nil
}

func GetCommentsByKeywords(keywords string, parentIDs []int) ([]models.HNComment, error) {
	var (
		comments []models.HNComment
		comment  models.HNComment
		err      error
		query    string
		args     []interface{}
	)
	params := map[string]interface{}{
		"keywords":   keywords,
		"parent_ids": parentIDs,
	}

	log.Info("keywords: ", keywords)
	log.Info("plain? :", isPlainTSQuery(keywords))
	if isPlainTSQuery(keywords) {
		query, args, err = sqlx.Named(stmtGetCommentsByKeywordsPlain, params)
	} else {
		query, args, err = sqlx.Named(stmtGetCommentsByKeywords, params)
	}

	if err != nil {
		return nil, err
	}

	query, args, err = sqlx.In(query, args...)
	if err != nil {
		return nil, err
	}

	query = db.Rebind(query)

	rows, err := db.Queryx(query, args...)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		err = rows.StructScan(&comment)
		if err != nil {
			log.Error(err)
		}
		comments = append(comments, comment)
	}

	return comments, nil
}
