package models

import (
	"time"
)

var Schema = ` 
CREATE TABLE IF NOT EXISTS hn_stories (
    id      integer         PRIMARY KEY     ,
    by      varchar(255)    NOT NULL        ,
    time    timestamp       NOT NULL        ,
    score   integer         NOT NULL        ,
    title   varchar(255)    NOT NULL        ,
    URL     varchar(255)    NOT NULL        
);

CREATE TABLE IF NOT EXISTS hn_comments (
    id          integer         PRIMARY KEY                         ,
    by          varchar(255)    NOT NULL                            ,
    time        timestamp       NOT NULL                            ,
    text        text            NOT NULL                            ,
    parent      integer         NOT NULL references hn_stories(id) 
);
`

var Indexes = `
DO $$
BEGIN

IF NOT EXISTS (
    SELECT 1
    FROM   pg_class c
    JOIN   pg_namespace n ON n.oid = c.relnamespace
    WHERE  c.relname = 'hn_comments_text_idx'
    AND    n.nspname = 'public'
    ) THEN

	CREATE INDEX hn_comments_text_idx ON hn_comments USING GIN(to_tsvector('english', text));
END IF;

END$$;
`

type HNPost struct {
	Id   int       `db:"id"    json:"id"`
	By   string    `db:"by"    json:"by"`
	Time time.Time `db:"time"  json:"time"`
}

type HNStory struct {
	HNPost
	Score int    `db:"score"   json:"score"`
	Title string `db:"title"   json:"title"`
	URL   string `db:"url"     json:"url"`
}

type HNComment struct {
	HNPost
	Text     string `db:"text"     json:"text"`
	ParentId int    `db:"parent"   json:"parent_id"`
}
