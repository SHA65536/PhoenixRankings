package main

import (
	"database/sql"
	"log"

	_ "embed"

	_ "github.com/mattn/go-sqlite3"
)

//go:embed queries/create_tables.sql
var create_tables string

//go:embed queries/create_indexes.sql
var create_indexes string

//go:embed queries/get_names.sql
var get_names string

type Database struct {
	DB          *sql.DB
	Names       map[string]int
	CurrentSnap *Snapshot
	LastSnap    *Snapshot
	Logger      *log.Logger
}

func MakeDatabase(path string, logger *log.Logger) (*Database, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}
	database := &Database{DB: db, Logger: logger}
	database.InitializeDB()

	return database, nil
}

// InitializeDB creates the database if none exists
// gets names map
func (db *Database) InitializeDB() error {
	var err error
	var name string
	var id int
	_, err = db.DB.Exec(create_tables)
	if err != nil {
		return err
	}
	_, err = db.DB.Exec(create_indexes)
	if err != nil {
		return err
	}
	rows, err := db.DB.Query(get_names)
	if err != nil {
		return err
	}
	defer rows.Close()

	db.Names = make(map[string]int)
	for rows.Next() {
		err = rows.Scan(&id, &name)
		if err != nil {
			return err
		}
		db.Names[name] = id
	}
	return nil
}
