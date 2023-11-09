package application

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/LNMMusic/tester/internal"
	"github.com/LNMMusic/tester/internal/cases"

	"github.com/go-sql-driver/mysql"
)

var (
	// ErrApplicationRun is the error of the run of the application.
	ErrApplicationRun = errors.New("application: run error")
	// ErrApplicationSetUp is the error of the set up of the application.
	ErrApplicationSetUp = errors.New("application: set up error")
	// ErrApplicationTearDown is the error of the tear down of the application.
	ErrApplicationTearDown = errors.New("application: tear down error")
)

// NewApplicationDefault creates a new application.
func NewApplicationDefault(cfg *Config) (a *ApplicationDefault) {
	// default config
	defaultConfig := &Config{
		Server: ServerConfig{
			Address: "http://localhost:8080",
		},
		Database: &DatabaseConfig{
			Address:  "localhost:3306",
		},
		Cases: CasesConfig{
			Reader: struct {
				FilePath  string
				BatchSize int
			}{
				FilePath:  "./cases.json",
				BatchSize: 10,
			},
		},
	}
	if cfg != nil {
		if cfg.Server.Address != "" {
			defaultConfig.Server.Address = cfg.Server.Address
		}
		if cfg.Database != nil {
			defaultConfig.Database = cfg.Database
		}
		if cfg.Cases.Reader.FilePath != "" {
			defaultConfig.Cases.Reader.FilePath = cfg.Cases.Reader.FilePath
		}
		if cfg.Cases.Reader.BatchSize > 0 {
			defaultConfig.Cases.Reader.BatchSize = cfg.Cases.Reader.BatchSize
		}
		if cfg.Cases.Reporter.ExcludedHeaders != nil {
			defaultConfig.Cases.Reporter.ExcludedHeaders = cfg.Cases.Reporter.ExcludedHeaders
		}
	}

	// application
	a = &ApplicationDefault{Config: defaultConfig}
	return
}

// Config is the config of the application.
type ServerConfig struct {
	// server address
	Address string
}
type DatabaseConfig struct {
	// database address
	Address string
	// database user
	User string
	// database password
	Password string
	// database name
	Name string
}
type CasesConfig struct {
	Reader struct {
		// cases file path
		FilePath string
		// batch size
		BatchSize int
	}
	Reporter struct {
		// excluded headers
		ExcludedHeaders []string
	}
}
// Config is the config of the application.
type Config struct {
	// server
	Server ServerConfig
	// database
	Database *DatabaseConfig
	// Cases
	Cases CasesConfig
}

// ApplicationDefault is the default implementation of Application.
// - cases: independent
// - database: centralized
// - http server: centralized
type ApplicationDefault struct {
	// configuration of the application
	*Config

	// instances of the application
	// - database
	db *sql.DB
	// - cases file
	casesFile *os.File
	// - reader
	reader *cases.ReaderJSON
	// - tester
	tester *internal.Tester
}

// TearDown tears down the application.
// - defer function, resources close in reverse order
func (a *ApplicationDefault) TearDown() (err error) {
	// close database
	err = a.db.Close()
	if err != nil {
		err = fmt.Errorf("%w. %v", ErrApplicationTearDown, err)
	}
	// close cases file
	e := a.casesFile.Close()
	if e != nil {
		err = fmt.Errorf("%w. %v. %w", ErrApplicationTearDown, e, err)
	}

	return
}

// SetUp sets up the application.
func (a *ApplicationDefault) SetUp() (err error) {
	// dependencies
	// - cases file
	a.casesFile, err = os.Open(a.Config.Cases.Reader.FilePath)
	if err != nil {
		return
	}
	// - database: connection
	cfg := &mysql.Config{
		Addr: a.Config.Database.Address,
		User: a.Config.Database.User,
		Passwd: a.Config.Database.Password,
		DBName: a.Config.Database.Name,
	}
	a.db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return
	}
	// - database: ping
	err = a.db.Ping()
	if err != nil {
		return
	}
	// - reader
	dc := json.NewDecoder(a.casesFile)
	ch := make(chan cases.CaseErr, a.Config.Cases.Reader.BatchSize)
	a.reader = cases.NewReaderJSON(dc, ch)
	// - tester
	ex := cases.NewDbExecuterMySQL(a.db)
	rq := cases.NewRequesterDefault(a.Config.Server.Address, &http.Client{})
	rp := cases.NewReporterDefault(a.Cases.Reporter.ExcludedHeaders)
	a.tester = internal.NewTester(ex, rq, rp)

	return
}

// Run runs the application.
func (a *ApplicationDefault) Run() (err error) {
	// stream cases
	go a.reader.Stream()

	// test cases
	for {
		// read case
		var c cases.Case
		c, err = a.reader.Read()
		if err != nil {
			if err == cases.ErrEndOfLine {
				err = nil
				break
			}
			return
		}

		// test case
		err = a.tester.Test(&c)
		if err != nil {
			return
		}
	}

	return
}