package main

import (
	"database/sql"
	"sync"
	"time"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

// Database to store a single connection to the database and its statistics
type Database struct {
	ID        uuid.UUID
	db        *sql.DB
	driver    string
	dsn       string
	query     string
	reconnect bool

	queries     float64 // Total number of queries
	queriesLock *sync.Mutex

	queriesDuration     time.Duration // Time spent executing queries
	queriesDurationLock *sync.Mutex
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
	database.ID = uuid.New()
	database.queriesLock = &sync.Mutex{}
	database.queriesDurationLock = &sync.Mutex{}
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
		queryStart := time.Now()
		result, err := d.db.Query(d.query)
		if err != nil {
			log.Fatalf("Failed query the database: %v", err)
		}
		if err = result.Close(); err != nil {
			log.Fatalf("Failed to close query: %v", err)
		}

		queryDuration := time.Now().Sub(queryStart)

		// Update statistics
		d.queriesLock.Lock()
		d.queries++
		d.queriesLock.Unlock()

		d.queriesDurationLock.Lock()
		d.queriesDuration = d.queriesDuration + queryDuration
		d.queriesDurationLock.Unlock()

		// Print debug statistics
		log.Debugf("Connection %s: Number of queries: %.f | Query duration: %s | Running sum of queries duration: %s", d.ID, d.queries, queryDuration, d.queriesDuration)

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

// AverageQueryTime returns the average query execution time in milliseconds
func (d *Database) AverageQueryTime() time.Duration {
	return d.queriesDuration / time.Duration(d.queries)
}
