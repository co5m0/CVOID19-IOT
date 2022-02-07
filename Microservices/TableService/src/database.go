package main

import (
	"database/sql"
	"errors"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type DoorMQ struct {
	Table     string `json:"table"`
	Seat      string `json:"seat"`
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
}

type database struct {
	filePath string
	conn     *sql.DB
}

func NewDatabase(dbpath string) (*database, error) {
	if _, err := os.Stat(dbpath); err != nil {
		log.Println("File not exists")
		return nil, errors.New("SQLite .db file passed does not exist")
	}
	conn, err := sql.Open("sqlite3", dbpath)
	if err != nil {
		return nil, err
	}
	db := database{filePath: dbpath, conn: conn}
	// status: 1 for open, 0 for close
	sqlStmt := `
	create table if not exists door (id integer not null primary key, table text, seat text, status text, timestamp text);
	`
	_, err = conn.Exec(sqlStmt)
	if err != nil {
		return nil, err
	}
	return &db, nil
}

func (db database) Close() {
	if db.conn != nil {
		db.conn.Close()
	}
}

func (db database) Insert(msgDoor DoorMQ) error {
	if db.conn == nil {
		return errors.New("cannot execute this operation, beacause the datebase connection is null")
	}
	stmt, err := db.conn.Prepare("INSERT INTO door(table, seat, status, timestamp) values(?,?,?,?)")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(msgDoor.Table, msgDoor.Seat, msgDoor.Status, msgDoor.Timestamp)
	if err != nil {
		return err
	}

	log.Println("Inset query executed")

	return nil
}

func (db database) GetAll() ([]DoorMQ, error) {
	if db.conn == nil {
		return nil, errors.New("cannot execute this operation, beacause the datebase connection is null")
	}
	var tableCount = 0
	doorsSlice := make([]DoorMQ, tableCount)
	var table string
	var seat string
	var status string
	var timestamp string

	// Count ho many rows are in the table
	err := db.conn.QueryRow("SELECT COUNT(*) FROM door").Scan(&tableCount)
	if err != nil {
		return nil, err
	}

	rows, err := db.conn.Query("SELECT table, seat, status, timestamp FROM door")
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		err = rows.Scan(&table, &seat, &status, &timestamp)
		if err != nil {
			return nil, err
		}
		doorsSlice = append(doorsSlice, DoorMQ{Table: table, Seat: seat, Status: status, Timestamp: timestamp})
		// log.Println("Row scanned", door, status, timestamp)
	}

	return doorsSlice, nil
}
