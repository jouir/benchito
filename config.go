package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"

	"gopkg.in/yaml.v2"
)

// Config to store all configurations (from command line, from file, etc)
type Config struct {
	Driver         string        `yaml:"driver"`
	Connections    int           `yaml:"connections"`
	Query          string        `yaml:"query"`
	Duration       time.Duration `yaml:"duration"`
	Reconnect      bool          `yaml:"reconnect"`
	DSN            string        `yaml:"dsn"`
	Host           string        `yaml:"host"`
	Port           int           `yaml:"port"`
	User           string        `yaml:"user"`
	Password       string        `yaml:"password"`
	Database       string        `yaml:"database"`
	TLS            string        `yaml:"tls"`
	ConnectTimeout int           `yaml:"connect_timeout"`
}

// NewConfig creates a Config struct
func NewConfig() *Config {
	return &Config{}
}

// Read YaML configuration file from disk
func (c *Config) Read(file string) error {
	file, err := filepath.Abs(file)
	if err != nil {
		return err
	}

	yamlFile, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		return err
	}

	return nil
}

// ParseDSN detects the database driver then creates the DSN accordingly
func (c *Config) ParseDSN() {
	if c.DSN == "" {
		switch c.Driver {
		case "postgres":
			c.DSN = c.parsePostgresDSN()
		case "mysql":
			c.DSN = c.parseMysqlDSN()
		}
	}
}

func (c *Config) parsePostgresDSN() string {
	var parameters []string
	if c.Host != "" {
		parameters = append(parameters, fmt.Sprintf("host=%s", c.Host))
	}
	if c.Port != 0 {
		parameters = append(parameters, fmt.Sprintf("port=%d", c.Port))
	}
	if c.User != "" {
		parameters = append(parameters, fmt.Sprintf("user=%s", c.User))
	}
	if c.Password != "" {
		parameters = append(parameters, fmt.Sprintf("password=%s", c.Password))
	}
	if c.Database != "" {
		parameters = append(parameters, fmt.Sprintf("database=%s", c.Database))
	}
	if c.ConnectTimeout != 0 {
		parameters = append(parameters, fmt.Sprintf("connect_timeout=%d", c.ConnectTimeout))
	}
	if AppName != "" {
		parameters = append(parameters, fmt.Sprintf("application_name=%s", AppName))
	}
	if c.TLS != "" {
		parameters = append(parameters, fmt.Sprintf("sslmode=%s", c.TLS))
	}
	return strings.Join(parameters, " ")
}

func (c *Config) parseMysqlDSN() (dsn string) {
	config := mysql.NewConfig()
	config.Timeout = 1 * time.Second
	config.Addr = c.Host
	if c.Port != 0 {
		config.Net = "tcp"
		config.Addr += fmt.Sprintf(":%d", c.Port)
	}
	config.User = c.User
	config.Passwd = c.Password
	config.DBName = c.Database
	config.TLSConfig = c.TLS
	return config.FormatDSN()
}
