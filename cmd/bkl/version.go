package main

import "fmt"

const version = "0.0.1"

var revision = ""

func getVersion() string {
	if revision == "" {
		return version
	}
	return fmt.Sprintf("%s (revision: %s)", version, revision)
}
