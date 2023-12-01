package main

import (
	"log"

	"github.com/underwoo16/gh-diffstack/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}

// For more examples of using go-gh, see:
// https://github.com/cli/go-gh/blob/trunk/example_gh_test.go
