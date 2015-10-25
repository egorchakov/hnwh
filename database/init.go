package database

import (
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/evgorchakov/hnwh/models"
)

const (
	dbEngine = "postgres"
)

var (
	db        *sqlx.DB
	log       = logrus.New()
	dbConnStr = os.Getenv("DATABASE_URL")
)

func SetupDB() {
	InitDB()
	db.MustExec(models.Schema)
	db.MustExec(models.Indexes)
}

func InitDB() {
	db = sqlx.MustConnect(dbEngine, dbConnStr)
}
