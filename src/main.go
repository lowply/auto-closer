package main

import (
	"log"
	"os"
)

func main() {
	// Required
	required := []string{"GITHUB_TOKEN", "GITHUB_REPOSITORY", "AC_LABEL"}
	for _, v := range required {
		if os.Getenv(v) == "" {
			log.Fatal(v + " is empty.")
		}
	}

	a, err := newAutoCloser()
	if err != nil {
		log.Fatal(err)
	}

	err = a.getIssuesList()
	if err != nil {
		log.Fatal(err)
	}

	err = a.closeIssues()
	if err != nil {
		log.Fatal(err)
	}
}
