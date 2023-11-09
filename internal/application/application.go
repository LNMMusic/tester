package application

// Application is an interface of an application.
type Application interface {
	// Run runs the application.
	Run() (err error)
}