package main

import (
	"fmt"

	"github.com/datawire/go-mkopensource/pkg/dependencies"
	"github.com/datawire/go-mkopensource/pkg/detectlicense"
	"github.com/datawire/go-mkopensource/pkg/golist"
	"github.com/datawire/go-mkopensource/pkg/util"
)

func GenerateDependencyList(
	modNames []string,
	modLicenses map[string]map[detectlicense.License]struct{},
	modInfos map[string]*golist.Module, goVersion string,
	licenseRestriction detectlicense.LicenseRestriction,
) (dependencyList dependencies.DependencyInfo, errors []error) {
	errors = []error{}

	for _, modKey := range modNames {
		ambassadorProprietary := isAmbassadorProprietary(modLicenses[modKey])
		if ambassadorProprietary {
			continue
		}

		modVal := modInfos[modKey]

		dependencyDetails := dependencies.Dependency{
			Name:     getDependencyName(modVal),
			Version:  getDependencyVersion(modVal, goVersion),
			Licenses: make(util.Set[detectlicense.License]),
		}
		for license := range modLicenses[modKey] {
			dependencyDetails.Licenses.Insert(license)
		}
		if errs := dependencyDetails.CheckLicenseRestrictions(licenseRestriction); errs != nil {
			errors = append(errors, errs...)
		}

		dependencyList.Dependencies = append(dependencyList.Dependencies, dependencyDetails)
	}

	return dependencyList, errors
}

func getDependencyName(modVal *golist.Module) string {
	if modVal == nil {
		return "the Go language standard library (\"std\")"
	}

	if modVal.Replace != nil && modVal.Replace.Version != "" && modVal.Replace.Path != modVal.Path {
		return fmt.Sprintf("%s (modified from %s)", modVal.Replace.Path, modVal.Path)
	}

	return modVal.Path
}

func getDependencyVersion(modVal *golist.Module, goVersion string) string {
	if modVal == nil {
		return goVersion
	}

	if modVal.Replace != nil {
		if modVal.Replace.Version == "" {
			return "(modified)"
		} else {
			return modVal.Replace.Version
		}
	}

	return modVal.Version
}
