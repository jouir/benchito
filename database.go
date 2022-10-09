package main

import (
	"database/sql"
	"log"
	"sync"
	"time"
)

// Database to store a single connection to the database and its statistics
type Database struct {
	db        *sql.DB
	driver    string
	dsn       string
	query     string
	reconnect bool
	queries   float64
}

// NewDatabase creates a single connection to the database then returns the struct
func NewDatabase(driver string, dsn string, query string, reconnect bool) (*Database, error) {
	database := &Database{
		driver:    driver,
		dsn:       dsn,
		query:     query,
		reconnect: reconnect,
	}
	err := database.connect()
	if err != nil {
		return nil, err
	}
	return database, nil
}

func (d *Database) connect() error {
	db, err := sql.Open(d.driver, d.dsn)
	if err != nil {
		return err
	}
	if err = db.Ping(); err != nil {
		return err
	}
	d.db = db
	return nil
}

func (d *Database) disconnect() error {
	if d.db != nil {
		return d.db.Close()
	}
	return nil
}

// Run to perform benchmark on a single connection for a duration
func (d *Database) Run(duration time.Duration, wg *sync.WaitGroup) error {
	defer wg.Done()
	defer d.disconnect()

	// Run single benchmark
	start := time.Now()
	end := start.Add(duration)
	for {
		if end.Before(time.Now()) {
			break
		}
		result, err := d.db.Query(d.query)
		if err != nil {
			log.Fatalf("Failed query the database: %v", err)
		}
		if err = result.Close(); err != nil {
			log.Fatalf("Failed to close query: %v", err)
		}

		d.queries++
		if d.reconnect {
			if err = d.disconnect(); err != nil {
				return err
			}
			if err = d.connect(); err != nil {
				return err
			}
		}

	}
	return nil
}

// Queries returns the number of performed queries in the benchmark
func (d *Database) Queries() float64 {
	return d.queries
}
