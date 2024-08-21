package main

import (
	"fmt"
	"os"

	"github.com/fe3dback/go-arch-lint-sdk/tests/projects/mvc/internal"
)

func main() {
	os.Exit(execute())
}

func execute() (exitCode int) {
	err := internal.Execute()
	if err != nil {
		fmt.Println(err)

		return 1
	}

	return 0
}
