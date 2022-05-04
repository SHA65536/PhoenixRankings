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

//go:embed queries/get_last_two.sql
var get_last_two string

type Database struct {
	DB          *sql.DB
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

// GetSnapshots fetches snapshots from db according
// to snapshot given
func (db *Database) GetSnapshots(reference *Snapshot) {
	var point *Datapoint
	db.Logger.Println("[Database] Making snapshots")
	db.CurrentSnap = &Snapshot{Players: make(map[int]*Datapoint)}
	db.LastSnap = &Snapshot{Players: make(map[int]*Datapoint)}
	for id := range reference.Players {
		point = &Datapoint{}
		rows, err := db.DB.Query(get_last_two, id)
		if err != nil {
			continue
		}
		if !rows.Next() {
			continue
		}
		rows.Scan(&point.DBId, &point.Id, &point.Name,
			&point.Rank, &point.Level, &point.Exp, &point.Fame,
			&point.Job, &point.Image, &point.Restriction,
		)
		db.CurrentSnap.Players[id] = point
		point = &Datapoint{}
		if !rows.Next() {
			continue
		}
		rows.Scan(&point.DBId, &point.Id, &point.Name,
			&point.Rank, &point.Level, &point.Exp, &point.Fame,
			&point.Job, &point.Image, &point.Restriction,
		)
		db.LastSnap.Players[id] = point
	}
	db.Logger.Printf("[Database] Fetched snapshots. Current: %d, Last: %d",
		len(db.CurrentSnap.Players), len(db.LastSnap.Players))
}

// InitializeDB creates the database if none exists
// gets names map
func (db *Database) InitializeDB() error {
	var err error
	_, err = db.DB.Exec(create_tables)
	if err != nil {
		return err
	}
	_, err = db.DB.Exec(create_indexes)
	if err != nil {
		return err
	}
	return nil
}

// SaveSnapshot recieves a snapshot to save to the database
// and commits the differences.
func (db *Database) SaveSnapshot(snap *Snapshot) {
	var changes int
	db.Logger.Println("[Database] Saving snapshot to database")
	for _, data := range snap.Players {
		if db.comparePoint(data) {
			changes++
			db.createPoint(data, snap.Timestamp)
		} else {
			db.updatePoint(data, snap.Timestamp)
		}
	}
	db.increaseSnap(snap)
	db.Logger.Printf("[Database] Snapshot saved. New points: %d", changes)
}

// comparePoint checks if the given datapoint
// needs creation, or just a change of timestamp.
func (db *Database) comparePoint(data *Datapoint) bool {
	cur, curOk := db.CurrentSnap.Players[data.Id]
	last, lastOk := db.LastSnap.Players[data.Id]
	if !(curOk && lastOk) {
		return true
	}
	if cur.Name != data.Name || last.Name != data.Name {
		return true
	}
	if cur.Rank != data.Rank || last.Rank != data.Rank {
		return true
	}
	if cur.Level != data.Level || last.Level != data.Level {
		return true
	}
	if cur.Exp != data.Exp || last.Exp != data.Exp {
		return true
	}
	if cur.Fame != data.Fame || last.Fame != data.Fame {
		return true
	}
	if cur.Job != data.Job || last.Job != data.Job {
		return true
	}
	return false
}

// createPoint creates a new data point in the database
func (db *Database) createPoint(data *Datapoint, timestamp int64) {
	//TODO: Update rowid
}

// updatePoint updates the timestamp of a datapoint
func (db *Database) updatePoint(data *Datapoint, timestamp int64) {
	//TODO: Update rowid
}

// increaseSnap updates the order of internal snapshots
func (db *Database) increaseSnap(snap *Snapshot) {
	db.LastSnap = db.CurrentSnap
	db.CurrentSnap = snap
}
