package main_test

import (
	main "github.com/datawire/go-mkopensource/cmd/go-mkopensource"
	. "github.com/datawire/go-mkopensource/pkg/detectlicense"
	"github.com/datawire/go-mkopensource/pkg/golist"
	"github.com/stretchr/testify/require"
	"testing"
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

	_, err := main.GenerateDependencyList(modNames, licenses, modInfos, goVersion, Unrestricted)
	require.NoError(t, err)

	_, err = main.GenerateDependencyList(modNames, licenses, modInfos, goVersion, OnAmbassadorServers)
	require.NoError(t, err)
}

func TestGenerateDependencyListWhenLicenseIsForbidden(t *testing.T) {
	licenses := map[string]map[License]struct{}{modNames[0]: {AGPL1Only: {}}}

	_, err := main.GenerateDependencyList(modNames, licenses, modInfos, goVersion, Unrestricted)
	require.Error(t, err)

	_, err = main.GenerateDependencyList(modNames, licenses, modInfos, goVersion, OnAmbassadorServers)
	require.Error(t, err)
}
