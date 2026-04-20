package version

import "fmt"

const version = "0.0.3"

var revision = ""

// Version returns the version and revision of the application.
func Version() string {
	if revision == "" {
		return version
	}
	return fmt.Sprintf("%s (revision: %s)", version, revision)
}
