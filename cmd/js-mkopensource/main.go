package main

import (
	"encoding/json"
	"fmt"
	"github.com/datawire/go-mkopensource/cmd/js-mkopensource/dependency"
	"github.com/datawire/go-mkopensource/pkg/detectlicense"
	"github.com/datawire/go-mkopensource/pkg/scanningerrors"
	"github.com/spf13/pflag"
	"os"
)

const (
	// Validations to do on the licenses.
	// The only validation for "internal" is to check chat forbidden licenses are not used
	internalApplication = "internal"
	// "external" applications have additional license requirements as documented in
	//https://www.notion.so/datawire/License-Management-5194ca50c9684ff4b301143806c92157
	externalApplication = "external"
)

type CLIArgs struct {
	ApplicationType string
}

func main() {
	args, err := parseArgs()
	if err != nil {
		if err == pflag.ErrHelp {
			os.Exit(int(NoError))
		}
		_, _ = fmt.Fprintf(os.Stderr, "%s: %v\nTry '%s --help' for more information.\n", os.Args[0], err, os.Args[0])
		os.Exit(int(InvalidArgumentsError))
	}

	licenseRestriction := getLicenseRestriction(args.ApplicationType)

	dependencyInfo, err := dependency.GetDependencyInformation(os.Stdin, licenseRestriction)
	if err != nil {
		err = scanningerrors.ExplainErrors([]error{err})
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

func parseArgs() (*CLIArgs, error) {
	args := &CLIArgs{}
	argparser := pflag.NewFlagSet(os.Args[0], pflag.ContinueOnError)
	help := false
	argparser.BoolVarP(&help, "help", "h", false, "Show this message")
	argparser.StringVar(&args.ApplicationType, "application-type", externalApplication,
		fmt.Sprintf("Where will the application run. One of: %s, %s\n"+
			"Internal applications are run on Ambassador servers.\n"+
			"External applications run on customer machines", internalApplication, externalApplication))

	if err := argparser.Parse(os.Args[1:]); err != nil {
		return nil, err
	}
	if help {
		fmt.Printf("Usage: %v OPTIONS\n", os.Args[0])
		fmt.Println()
		fmt.Println("OPTIONS:")
		argparser.PrintDefaults()
		return nil, pflag.ErrHelp
	}

	if argparser.NArg() != 0 {
		return nil, fmt.Errorf("expected 0 arguments, got %d: %q", argparser.NArg(), argparser.Args())
	}

	if args.ApplicationType != internalApplication && args.ApplicationType != externalApplication {
		return nil, fmt.Errorf("--application-type must be one of '%s', '%s'", internalApplication, externalApplication)
	}

	return args, nil
}

func getLicenseRestriction(applicationType string) detectlicense.LicenseRestriction {
	var LicenseRestriction detectlicense.LicenseRestriction
	switch applicationType {
	case internalApplication:
		LicenseRestriction = detectlicense.AmbassadorServers
	default:
		LicenseRestriction = detectlicense.Unrestricted
	}
	return LicenseRestriction
}
