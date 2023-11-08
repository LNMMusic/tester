package cases

import (
	"database/sql"
	"fmt"
)

// NewDbExecuterMySQL creates a new MySQL database executer.
func NewDbExecuterMySQL(db *sql.DB) *DbExecuterMySQL {
	return &DbExecuterMySQL{
		db: db,
	}
}

// DbExecuterMySQL is a MySQL database executer.
type DbExecuterMySQL struct {
	// db is the database to execute queries on.
	db *sql.DB
}

// Exec executes queries on the database.
func (e *DbExecuterMySQL) Exec(queries ...string) (err error) {
	// execute queries
	for i, q := range queries {
		_, err = e.db.Exec(q)
		if err != nil {
			err = fmt.Errorf("error executing query %d", i)
			return
		}
	}

	return
}