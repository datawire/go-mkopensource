package main_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	main "github.com/datawire/go-mkopensource/cmd/go-mkopensource"
	. "github.com/datawire/go-mkopensource/pkg/detectlicense"
	"github.com/datawire/go-mkopensource/pkg/golist"
)

//nolint:gochecknoglobals // Can't be a constant
var (
	goVersion = "v1.17.3"
	modNames  = []string{"github.com/josharian/intern"}
	modInfos  = map[string]*golist.Module{
		"github.com/josharian/intern": {
			Path:    "github.com/josharian/intern",
			Version: "1.2.3",
		}}
)

func TestGenerateDependencyListWhenLicenseIsAllowed(t *testing.T) {
	licenses := map[string]map[License]struct{}{modNames[0]: {BSD1: {}}}

	_, errors := main.GenerateDependencyList(modNames, licenses, modInfos, goVersion, Unrestricted)
	require.Empty(t, errors)

	_, errors = main.GenerateDependencyList(modNames, licenses, modInfos, goVersion, AmbassadorServers)
	require.Empty(t, errors)
}

func TestGenerateDependencyListWhenLicenseIsForbidden(t *testing.T) {
	licenses := map[string]map[License]struct{}{modNames[0]: {AGPL1Only: {}}}

	_, errors := main.GenerateDependencyList(modNames, licenses, modInfos, goVersion, Unrestricted)
	require.NotEmptyf(t, errors, "Expected at least one error but got none")

	_, errors = main.GenerateDependencyList(modNames, licenses, modInfos, goVersion, AmbassadorServers)
	require.NotEmptyf(t, errors, "Expected at least one error but got none")
}
