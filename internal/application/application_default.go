package application

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/LNMMusic/tester/internal"
	"github.com/LNMMusic/tester/internal/cases"
	"github.com/go-sql-driver/mysql"
)

var (
	// ErrApplicationRun is the error of the run of the application.
	ErrApplicationRun = errors.New("application: run error")
)

// NewApplicationDefault creates a new application.
func NewApplicationDefault(cfg *Config) (a *ApplicationDefault) {
	// default config
	defaultCfg := &Config{
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
		defaultCfg = cfg
	}

	// application
	a = &ApplicationDefault{
		cfg: defaultCfg,
	}
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
type ApplicationDefault struct {
	// configuration of the application
	cfg *Config
}

// Run runs the application.
func (a *ApplicationDefault) Run() (err error) {
	// dependency injection
	// - reader: file
	f, err := os.Open(a.cfg.Cases.Reader.FilePath)
	if err != nil {
		err = fmt.Errorf("%w - %v", ErrApplicationRun, err)
		return
	}
	defer f.Close()
	dc := json.NewDecoder(f)
	// - reader: chan
	ch := make(chan cases.CaseErr, a.cfg.Cases.Reader.BatchSize)
	rd := cases.NewReaderJSON(dc, ch)

	// - casetester: dbexecuter
	cfg := &mysql.Config{
		User:                 a.cfg.Database.User,
		Passwd:               a.cfg.Database.Password,
		Net:                  "tcp",
		Addr:                 a.cfg.Database.Address,
		DBName:               a.cfg.Database.Name,
	}
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		err = fmt.Errorf("%w - %v", ErrApplicationRun, err)
		return
	}
	defer db.Close()
	ex := cases.NewDbExecuterMySQL(db)
	// - casetester: requester
	rq := cases.NewRequesterDefault(a.cfg.Server.Address, nil)
	// - casetester: reporter
	rp := cases.NewReporterDefault(a.cfg.Cases.Reporter.ExcludedHeaders)
	// - casetester: case tester
	ct := internal.NewCaseTesterDefault(ex, rq, rp)

	// - tester
	ts := internal.NewTester(rd, ct)

	// run
	// - stream cases
	go rd.Stream()
	// - test cases
	err = ts.Run()
	if err != nil {
		err = fmt.Errorf("%w - %v", ErrApplicationRun, err)
		return
	}

	return
}