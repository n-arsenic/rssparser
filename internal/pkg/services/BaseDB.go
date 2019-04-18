package services

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"rssparser/internal/pkg/config"
)

type Postgres struct {
	Db *sql.DB
}

var Postgre *Postgres = &Postgres{
	&sql.DB{},
}

func (postgre *Postgres) Connect() {
	var conf *config.Config = config.New()

	db, err := sql.Open("postgres", fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s sslmode=disable",
		conf.DB_USER, conf.DB_PASSWORD, conf.DB_NAME, conf.DB_HOST))
	if err != nil {
		log.Fatal(err)
	}
	postgre.Db = db
}

func (postgre *Postgres) Close() (err error) {
	if postgre.Db == nil {
		return
	}
	err = postgre.Db.Close()
	return err
}
