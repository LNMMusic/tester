package cases

// DbExecuter is an interface for executing queries on a database.
type DbExecuter interface {
	// Exec executes queries on the database.
	Exec(queries ...string) (err error)
}