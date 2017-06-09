package registry

type Backend interface {
	// Register registers fabio as a service in the registry.
	Register() error

	// Deregister removes the service registration for fabio.
	Deregister() error

	// ReadManual returns the current manual overrides and
	// their version as seen by the registry.
	ReadManual() (value string, version uint64, err error)

	// WriteManual writes the new value to the registry if the
	// version of the stored document still matchhes version.
	WriteManual(value string, version uint64) (ok bool, err error)

	// WatchServices watches the registry for changes in service
	// registration and health and pushes them if there is a difference.
	WatchServices() chan string

	// WatchManual watches the registry for changes in the manual
	// overrides and pushes them if there is a difference.
	WatchManual() chan string
}

var Default Backend
