package dbutils

import (
	"github.com/evgorchakov/hnwh/models"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"strings"
)

const (
	dBEngine  = "postgres"
	dBConnStr = "user=dev dbname=hnwh sslmode=disable"
)

var (
	queries = map[string]string{
		"insert_or_update_story": `
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
		`,

		"insert_or_update_comment": `
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
		`,
	}
)

func OpenDB() *sqlx.DB {
	db, _ := sqlx.Open(dBEngine, dBConnStr)
	return db
}

func InitDB() {
	db := sqlx.MustConnect(dBEngine, dBConnStr)
	defer db.Close()
	db.MustExec(models.Schema)
	db.MustExec(models.Indexes)
}

func InsertOrUpdateStory(story *models.HNStory) error {
	db := OpenDB()
	defer db.Close()
	_, err := db.NamedExec(queries["insert_or_update_story"], story)
	if err != nil {
		log.Println(err)
	}
	return err
}

func InsertOrUpdateComment(comment *models.HNComment) error {
	db := OpenDB()
	defer db.Close()
	_, err := db.NamedExec(queries["insert_or_update_comment"], comment)
	if err != nil {
		log.Println(err)
	}
	return err
}

func BuildFullTextSearhQuery(params map[string]interface{}) string {
	var final_query, select_query, anded_where_clauses string
	var where_clauses []string

	select_query = `SELECT * FROM hn_comments`

	query_elements := map[string]string{
		"tags": `text @@ to_tsquery(:tags)`,
		"from": `time >= :from`,
		"to":   `time <= :to`,
	}

	for k, v := range query_elements {
		if _, ok := params[k]; ok {
			where_clauses = append(where_clauses, v)
		}
	}

	anded_where_clauses = strings.Join(where_clauses, " AND ")
	final_query = select_query + " WHERE " + anded_where_clauses
	// final_query = strings.Join([]string{select_query, "WHERE", anded_where_clauses}, " ")

	return final_query
}
