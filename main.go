package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

// AppName to store application name
var AppName string = "benchito"

// AppVersion to set version at compilation time
var AppVersion string = "9999"

// GitCommit to set git commit at compilation time (can be empty)
var GitCommit string

// GoVersion to set Go version at compilation time
var GoVersion string

func main() {

	config := NewConfig()

	version := flag.Bool("version", false, "Print version and exit")
	configFile := flag.String("config", "", "Configuration file")
	flag.StringVar(&config.Driver, "driver", "postgres", "Database driver (postgres or mysql)")
	flag.IntVar(&config.Connections, "connections", 1, "Number of concurrent connections to the database")
	flag.StringVar(&config.Query, "query", "SELECT /* "+AppName+" */ NOW();", "Query to execute for the benchmark")
	flag.DurationVar(&config.Duration, "duration", 1*time.Second, "Duration of the benchmark")
	flag.BoolVar(&config.Reconnect, "reconnect", false, "Force database reconnection between each queries")
	flag.StringVar(&config.DSN, "dsn", "", "Database cpnnection string")
	flag.StringVar(&config.Host, "host", "", "Host address of the database")
	flag.IntVar(&config.Port, "port", 0, "Port of the database")
	flag.StringVar(&config.User, "user", "", "Username of the database")
	flag.StringVar(&config.Password, "password", "", "Password of the database")
	flag.StringVar(&config.Database, "database", "", "Database name")
	flag.StringVar(&config.TLS, "tls", "", "TLS configuration")
	flag.Parse()

	if *version {
		showVersion()
		return
	}

	if *configFile != "" {
		err := config.Read(*configFile)
		if err != nil {
			log.Fatalf("Failed to read configuration file: %v", err)
		}
	}

	config.ParseDSN()

	benchmark, err := NewBenchmark(config.Connections, config.Duration, config.Driver, config.DSN, config.Query, config.Reconnect)
	if err != nil {
		log.Fatalf("Cannot perform benchmark: %v", err)
	}
	benchmark.Run()
	log.Printf("Queries: %.0f", benchmark.Queries())
	log.Printf("Queries per second: %.0f", benchmark.QueriesPerSecond())
}

func showVersion() {
	if GitCommit != "" {
		AppVersion = fmt.Sprintf("%s-%s", AppVersion, GitCommit)
	}
	fmt.Printf("%s version %s (compiled with %s)\n", AppName, AppVersion, GoVersion)
}
