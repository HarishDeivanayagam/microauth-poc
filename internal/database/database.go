package database

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Database struct {
	client *sqlx.DB
}

func New() *Database {
	return &Database{}
}

func (d *Database) Connect(driver string, url string) error {
	db, err := sqlx.Connect(driver, url)
	if err != nil {
		return err
	}
	d.client = db
	return nil
}
