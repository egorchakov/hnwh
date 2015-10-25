package database

import (
	"github.com/evgorchakov/hnwh/models"
)

const (
	stmtInsertOrUpdateStory = `
		WITH new_values (id, by, time, score, title, url) as 
			(values (:id, :by, :time, :score, :title, :url)),

		upsert as (
			update hn_stories m
				set by = nv.by,
					time = CAST(nv.time AS timestamp),
					score = CAST(nv.score AS integer), 
					title = nv.title,
					url = nv.url
			FROM new_values nv
			WHERE m.id = CAST(nv.id AS integer)
			RETURNING m.*
		)

		INSERT INTO hn_stories (id, by, time, score, title, url)
		SELECT CAST(id AS integer), by, CAST(time AS timestamp), CAST(score AS integer), title, url
		FROM new_values
		WHERE NOT EXISTS (SELECT 1 from upsert up WHERE up.id = CAST(new_values.id AS integer) )
	`

	stmtInsertOrUpdateComment = `
		WITH new_values (id, by, time, text, parent) as 
			(values (:id, :by, :time, :text, :parent)),

		upsert as (
			update hn_comments m
				set by = nv.by,
					time = CAST(nv.time AS timestamp),
					text = nv.text,
					parent = CAST(nv.parent AS integer)
			FROM new_values nv
			WHERE m.id = CAST(nv.id AS integer)
			RETURNING m.*
		)

		INSERT INTO hn_comments (id, by, time, text, parent)
		SELECT CAST(id AS integer), by, CAST(time AS timestamp), text, CAST(parent AS integer)
		FROM new_values
		WHERE NOT EXISTS (SELECT 1 from upsert up WHERE up.id = CAST(new_values.id AS integer))
	`
)

func InsertOrUpdateStory(story *models.HNStory) error {
	_, err := db.NamedExec(stmtInsertOrUpdateStory, story)
	return err
}

func InsertOrUpdateComment(comment *models.HNComment) error {
	_, err := db.NamedExec(stmtInsertOrUpdateComment, comment)
	return err
}
