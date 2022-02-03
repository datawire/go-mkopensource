package main

import (
	"encoding/json"
	"fmt"
	"github.com/datawire/go-mkopensource/cmd/js-mkopensource/dependency"
	"os"
)

func main() {
	dependencyInfo, err := dependency.GetDependencyInformation(os.Stdin)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error generating dependency information: %v\n", err)
		os.Exit(int(DependencyGenerationError))
	}

	jsonString, marshalErr := json.Marshal(dependencyInfo)
	if marshalErr != nil {
		_, _ = fmt.Fprintf(os.Stderr, "could not generate JSON output: %v\n", err)
		os.Exit(int(MarshallJsonError))
	}

	if _, err := os.Stdout.Write(jsonString); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "could not write JSON output: %v\n", err)
		os.Exit(int(WriteError))
	}

	_, _ = fmt.Fprintf(os.Stdout, "\n")
}