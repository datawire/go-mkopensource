package dependency

import (
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/datawire/go-mkopensource/pkg/dependencies"
	"github.com/datawire/go-mkopensource/pkg/detectlicense"
	"github.com/datawire/go-mkopensource/pkg/scanningerrors"
	"github.com/datawire/go-mkopensource/pkg/util"
)

type NodeDependencies map[string]nodeDependency

type nodeDependency struct {
	Licenses       interface{} `json:"licenses"`
	Repository     string      `json:"repository"`
	DependencyPath string      `json:"dependencyPath"`
	Name           string      `json:"name"`
	Version        string      `json:"version"`
	Path           string      `json:"path"`
	URL            string      `json:"url"`
	LicenseFile    string      `json:"licenseFile"`
	LicenseText    string      `json:"licenseText"`
}

func (n *nodeDependency) licenses() (string, error) {
	if licenses, ok := n.Licenses.(string); ok {
		return licenses, nil
	}

	if licenseArray, ok := n.Licenses.([]interface{}); ok {
		var licenses []string
		for _, v := range licenseArray {
			license, ok := v.(string)
			if !ok {
				return "", fmt.Errorf("Dependency '%s@%s' has an invalid license field: %#v", n.Name, n.Version, n.Licenses)
			}
			licenses = append(licenses, license)
		}

		return strings.Join(licenses, " AND "), nil
	}

	return "", fmt.Errorf("Dependency '%s@%s' has an invalid license field: %v", n.Name, n.Version, n.Licenses)
}

func GetDependencyInformation(r io.Reader, licenseRestriction detectlicense.LicenseRestriction) (dependencyInfo dependencies.DependencyInfo, err error) {
	var nodeDependencies NodeDependencies
	data, err := io.ReadAll(r)
	if err != nil {
		return
	}

	if err := json.Unmarshal(data, &nodeDependencies); err != nil {
		return dependencies.DependencyInfo{}, err
	}

	licErrs := []error{}
	for _, dependencyId := range util.SortedMapKeys(nodeDependencies) {
		nodeDependency := nodeDependencies[dependencyId]

		dependency, dependencyErr := getDependencyDetails(nodeDependency, dependencyId)
		if dependencyErr != nil {
			licErrs = append(licErrs, dependencyErr)
			continue
		}

		if licenseErr := dependency.CheckLicenseRestrictions(licenseRestriction); licenseErr != nil {
			licErrs = append(licErrs, licenseErr...)
			continue
		}

		dependencyInfo.Dependencies = append(dependencyInfo.Dependencies, *dependency)
	}

	if len(licErrs) > 0 {
		return dependencyInfo, scanningerrors.ExplainErrors(licErrs)
	}

	return dependencyInfo, err
}

func getDependencyDetails(nodeDependency nodeDependency, dependencyId string) (*dependencies.Dependency, error) {
	name, version := splitDependencyIdentifier(dependencyId)

	dependency := &dependencies.Dependency{
		Name:    name,
		Version: version,
	}

	allLicenses, err := getDependencyLicenses(dependencyId, nodeDependency)
	if err != nil {
		return nil, err
	}
	dependency.Licenses = allLicenses

	return dependency, nil
}

func getDependencyLicenses(dependencyId string, nodeDependency nodeDependency) (util.Set[detectlicense.License], error) {
	licenseString, err := nodeDependency.licenses()
	if err != nil {
		return nil, err
	}

	if licenseString == "" {
		return nil, fmt.Errorf("Dependency '%s@%s' is missing a license identifier.", nodeDependency.Name, nodeDependency.Version)
	}

	parenthesisRe, err := regexp.Compile(`^\(|\)$`)
	if err != nil {
		return nil, err
	}
	licenseString = parenthesisRe.ReplaceAllString(licenseString, "")

	separatorRe, err := regexp.Compile(` OR | AND `)
	if err != nil {
		return nil, err
	}
	licenses := separatorRe.Split(licenseString, -1)

	allLicenses := make(util.Set[detectlicense.License])
	for _, spdxId := range licenses {
		if license, ok := detectlicense.SpdxIdentifiers[spdxId]; ok {
			allLicenses.Insert(license)
			continue
		}

		if licenses, ok := hardcodedJsDependencies[dependencyId]; ok {
			for _, lic := range licenses {
				allLicenses.Insert(lic)
			}
			break
		}

		return nil, fmt.Errorf("Dependency '%s@%s' has an unknown SPDX Identifier '%s'.",
			nodeDependency.Name, nodeDependency.Version, spdxId)
	}

	return allLicenses, nil
}

func splitDependencyIdentifier(identifier string) (name string, version string) {
	parts := strings.Split(identifier, "@")

	numberOfParts := len(parts)
	return strings.Join(parts[:numberOfParts-1], "@"), parts[numberOfParts-1]
}
