package main

import "fmt"

const version = "0.0.2"

var revision = ""

func getVersion() string {
	if revision == "" {
		return version
	}
	return fmt.Sprintf("%s (revision: %s)", version, revision)
}
