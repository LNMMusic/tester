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

// ConfigApplicationDefault is the config of the default application.
type ConfigApplicationDefault struct {
	// cases file path
	CasesFilePath string
	// database config
	Database *mysql.Config
	// http server address
	ServerAddress string
}

// NewApplicationDefault creates a new application.
func NewApplicationDefault(cfg *ConfigApplicationDefault) (a *ApplicationDefault) {
	// default config
	defaultCfg := &ConfigApplicationDefault{
		CasesFilePath: "cases.json",
		Database: &mysql.Config{},
		ServerAddress: "http://localhost:8080",
	}
	if cfg != nil {
		if cfg.CasesFilePath != "" {
			defaultCfg.CasesFilePath = cfg.CasesFilePath
		}
		if cfg.Database != nil {
			defaultCfg.Database = cfg.Database
		}
		if cfg.ServerAddress != "" {
			defaultCfg.ServerAddress = cfg.ServerAddress
		}
	}

	// application
	a = &ApplicationDefault{
		casesFilePath: defaultCfg.CasesFilePath,
		database: defaultCfg.Database,
		serverAddress: defaultCfg.ServerAddress,
	}
	return
}


// ApplicationDefault is the default implementation of Application.
// - cases: independent
// - database: centralized
// - http server: centralized
type ApplicationDefault struct {
	// cases file
	casesFilePath string
	casesFile *os.File
	
	// database config
	database *mysql.Config
	db *sql.DB
	
	// http server config
	serverAddress string

	// reader is the reader of cases.
	reader *cases.ReaderJSON
	
	// tester is the tester controller.
	tester internal.Tester
}

// TearDown tears down the application.
// - defer function, resources close in reverse order
func (a *ApplicationDefault) TearDown() (err error) {
	defer func() {
		e := a.casesFile.Close()
		if e != nil {
			err = fmt.Errorf("%w. %v. %w", ErrApplicationTearDown, e, err)
		}
	}()

	defer func() {
		e := a.db.Close()
		if e != nil {
			err = fmt.Errorf("%w. %v. %w", ErrApplicationTearDown, e, err)
		}
	}()
	// // close database
	// e := a.db.Close()
	// if e != nil {
	// 	err = fmt.Errorf("%w. %v. %w", ErrApplicationTearDown, e, err)
	// }
	// // close cases file
	// e = a.casesFile.Close()
	// if e != nil {
	// 	err = fmt.Errorf("%w. %v. %w", ErrApplicationTearDown, e, err)
	// }
	return
}


// SetUp sets up the application.
func (a *ApplicationDefault) SetUp() (err error) {
	// dependencies
	// - cases file
	a.casesFile, err = os.Open(a.casesFilePath)
	if err != nil {
		return
	}
	// - database
	a.db, err = sql.Open("mysql", a.database.FormatDSN())
	if err != nil {
		return
	}
	
	// - reader
	dc := json.NewDecoder(a.casesFile)
	ch := make(chan cases.CaseErr, 10)
	a.reader = cases.NewReaderJSON(dc, ch)
	// - tester
	ex := cases.NewDbExecuterMySQL(a.db)
	rq := cases.NewRequesterDefault(a.serverAddress, &http.Client{
		Timeout: 5,
	})
	rp := cases.NewReporterDefault()
	a.tester = *internal.NewTester(ex, rq, rp)

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