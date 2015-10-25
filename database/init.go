package database

import (
	"github.com/Sirupsen/logrus"
	"github.com/evgorchakov/hnwh/models"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const (
	dbEngine  = "postgres"
	dbConnStr = "user=dev dbname=hnwh sslmode=disable"
)

var (
	db  *sqlx.DB
	log = logrus.New()
)

func SetupDB() {
	InitDB()
	db.MustExec(models.Schema)
	db.MustExec(models.Indexes)
}

func InitDB() {
	db = sqlx.MustConnect(dbEngine, dbConnStr)
}
