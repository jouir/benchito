package main

import (
	"sync"
	"time"
)

// Benchmark to store benchmark state
type Benchmark struct {
	duration    time.Duration
	connections int
	databases   []*Database
}

// NewBenchmark connects to the database then creates a Benchmark struct
func NewBenchmark(connections int, duration time.Duration, driver string, dsn string, query string, reconnect bool) (*Benchmark, error) {
	var databases []*Database
	for i := 0; i < connections; i++ {
		database, err := NewDatabase(driver, dsn, query, reconnect)
		if err != nil {
			return nil, err
		}
		databases = append(databases, database)
	}
	return &Benchmark{
		duration:    duration,
		connections: connections,
		databases:   databases,
	}, nil
}

// Run performs the benchmark by runing queries for a duration
func (b *Benchmark) Run() {
	wg := new(sync.WaitGroup)
	wg.Add(b.connections)
	for _, database := range b.databases {
		go database.Run(b.duration, wg)
	}
	wg.Wait()
}

// Queries returns the number of executed queries during the benchmark
func (b *Benchmark) Queries() (queries float64) {
	for _, database := range b.databases {
		queries = queries + database.Queries()
	}
	return
}

// QueriesPerSecond returns the number of executed queries per second during the benchmark
func (b *Benchmark) QueriesPerSecond() float64 {
	return b.Queries() / b.duration.Seconds()
}
