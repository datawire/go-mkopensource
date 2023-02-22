package dependencies_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/datawire/go-mkopensource/pkg/dependencies"
	"github.com/datawire/go-mkopensource/pkg/detectlicense"
	"github.com/datawire/go-mkopensource/pkg/util"
)

//nolint:gochecknoglobals // Can't be a constant
var (
	emptyDependencies = dependencies.DependencyInfo{
		Dependencies: []dependencies.Dependency{},
	}

	dependenciesWithUniqueLicenses = dependencies.DependencyInfo{
		Dependencies: []dependencies.Dependency{
			{
				Name:     "library1",
				Version:  "1.0.2",
				Licenses: util.NewSet(detectlicense.MIT),
			},
			{
				Name:     "library2",
				Version:  "3.1.2",
				Licenses: util.NewSet(detectlicense.BSD1),
			},
		},
	}

	dependencyWithMultipleLicenses = dependencies.DependencyInfo{
		Dependencies: []dependencies.Dependency{
			{
				Name:     "library1",
				Version:  "1.0.2",
				Licenses: util.NewSet(detectlicense.GPL3Only, detectlicense.BSD2),
			},
		},
	}

	dependenciesWithOverlappingLicenses = dependencies.DependencyInfo{
		Dependencies: []dependencies.Dependency{
			{
				Name:     "library1",
				Version:  "1.0.2",
				Licenses: util.NewSet(detectlicense.GPL3Only, detectlicense.BSD2),
			},
			{
				Name:     "library2",
				Version:  "3.1.2",
				Licenses: util.NewSet(detectlicense.BSD2),
			},
			{
				Name:     "library2",
				Version:  "3.1.2",
				Licenses: util.NewSet(detectlicense.Apache2, detectlicense.GPL3Only),
			},
		},
	}

	licensesWithoutUrls = dependencies.DependencyInfo{
		Dependencies: []dependencies.Dependency{
			{
				Name:     "library1",
				Version:  "1.0.2",
				Licenses: util.NewSet(detectlicense.PublicDomain),
			},
		},
	}
)

func TestLicenseListIsCorrect(t *testing.T) {
	testCases := []struct {
		Name             string
		dependencies     dependencies.DependencyInfo
		expectedLicenses map[string]string
	}{
		{
			"Empty dependency list",
			emptyDependencies,
			map[string]string{},
		},
		{
			"Several dependencies with different licenses",
			dependenciesWithUniqueLicenses,
			map[string]string{
				detectlicense.MIT.Name:  detectlicense.MIT.URL,
				detectlicense.BSD1.Name: detectlicense.BSD1.URL},
		},
		{
			"A dependency with multiple licenses",
			dependencyWithMultipleLicenses,
			map[string]string{
				detectlicense.GPL3Only.Name: detectlicense.GPL3Only.URL,
				detectlicense.BSD2.Name:     detectlicense.BSD2.URL,
			},
		},
		{
			"Dependencies with overlapping licenses",
			dependenciesWithOverlappingLicenses,
			map[string]string{
				detectlicense.GPL3Only.Name: detectlicense.GPL3Only.URL,
				detectlicense.BSD2.Name:     detectlicense.BSD2.URL,
				detectlicense.Apache2.Name:  detectlicense.Apache2.URL,
				detectlicense.GPL3Only.Name: detectlicense.GPL3Only.URL,
			},
		},
		{
			"Licenses without Url",
			licensesWithoutUrls,
			map[string]string{
				detectlicense.PublicDomain.Name: "",
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			var jsonObj struct {
				Dependencies []dependencies.Dependency `json:"dependencies"`
				Licenses     map[string]string         `json:"licenseInfo"`
			}
			jsonBytes, err := json.Marshal(testCase.dependencies)
			require.NoError(t, err)
			require.NoError(t, json.Unmarshal(jsonBytes, &jsonObj))

			require.Equal(t, testCase.expectedLicenses, jsonObj.Licenses)
		})
	}
}

func TestCheckLicensesValidatesAllowedLicenseCorrectly(t *testing.T) {
	testCases := []struct {
		testName           string
		licenses           util.Set[detectlicense.License]
		licenseRestriction detectlicense.LicenseRestriction
	}{
		{
			"Unrestricted licenses are OK on Ambassador Labs servers",
			util.NewSet(detectlicense.MIT),
			detectlicense.AmbassadorServers,
		},
		{
			"Unrestricted licenses are OK everywhere",
			util.NewSet(detectlicense.MIT),
			detectlicense.Unrestricted,
		},
		{
			"Restricted licenses are OK on Ambassador Labs servers",
			util.NewSet(detectlicense.GPL3Only),
			detectlicense.AmbassadorServers,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			testDependency := dependencies.Dependency{
				Name:     "library1",
				Version:  "1.0.2",
				Licenses: testCase.licenses,
			}
			errs := testDependency.CheckLicenseRestrictions(testCase.licenseRestriction)
			require.Len(t, errs, 0)
		})
	}
}

func TestCheckLicensesValidatesForbiddenLicensesCorrectly(t *testing.T) {
	testCases := []struct {
		testName           string
		licenses           util.Set[detectlicense.License]
		licenseRestriction detectlicense.LicenseRestriction
	}{
		{
			"It's not possible to allow the use of forbidden licenses by mistake",
			util.NewSet(detectlicense.AGPL1Only),
			detectlicense.Forbidden,
		},
		{
			"Forbidden licenses are not allowed on Ambassador Labs servers",
			util.NewSet(detectlicense.AGPL1Only),
			detectlicense.AmbassadorServers,
		},
		{
			"Forbidden licenses are not allowed on customer machines",
			util.NewSet(detectlicense.AGPL1Only),
			detectlicense.Unrestricted,
		},
		{
			"Restricted licenses are not OK on customer machines",
			util.NewSet(detectlicense.GPL3Only),
			detectlicense.Unrestricted,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			testDependency := dependencies.Dependency{
				Name:     "library1",
				Version:  "1.0.2",
				Licenses: testCase.licenses,
			}
			errs := testDependency.CheckLicenseRestrictions(testCase.licenseRestriction)
			require.Len(t, errs, 1)
		})
	}
}
