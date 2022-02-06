package dependencies_test

import (
	"github.com/datawire/go-mkopensource/pkg/dependencies"
	"github.com/datawire/go-mkopensource/pkg/detectlicense"
	"github.com/stretchr/testify/require"
	"testing"
)

//nolint:gochecknoglobals // Can't be a constant
var (
	emptyDependencies = dependencies.DependencyInfo{
		Dependencies: []dependencies.Dependency{},
		Licenses:     map[string]string{},
	}

	dependenciesWithUniqueLicenses = dependencies.DependencyInfo{
		Dependencies: []dependencies.Dependency{
			{
				Name:     "library1",
				Version:  "1.0.2",
				Licenses: []string{detectlicense.MIT.Name},
			},
			{
				Name:     "library2",
				Version:  "3.1.2",
				Licenses: []string{detectlicense.BSD1.Name},
			},
		},
		Licenses: map[string]string{},
	}

	dependencyWithMultipleLicenses = dependencies.DependencyInfo{
		Dependencies: []dependencies.Dependency{
			{
				Name:     "library1",
				Version:  "1.0.2",
				Licenses: []string{detectlicense.GPL3Only.Name, detectlicense.BSD2.Name},
			},
		},
		Licenses: map[string]string{},
	}

	dependenciesWithOverlappingLicenses = dependencies.DependencyInfo{
		Dependencies: []dependencies.Dependency{
			{
				Name:     "library1",
				Version:  "1.0.2",
				Licenses: []string{detectlicense.GPL3Only.Name, detectlicense.BSD2.Name},
			},
			{
				Name:     "library2",
				Version:  "3.1.2",
				Licenses: []string{detectlicense.BSD2.Name},
			},
			{
				Name:     "library2",
				Version:  "3.1.2",
				Licenses: []string{detectlicense.Apache2.Name, detectlicense.GPL3Only.Name},
			},
		},
		Licenses: map[string]string{},
	}

	licensesWithoutUrls = dependencies.DependencyInfo{
		Dependencies: []dependencies.Dependency{
			{
				Name:     "library1",
				Version:  "1.0.2",
				Licenses: []string{detectlicense.PublicDomain.Name},
			},
		},
		Licenses: map[string]string{},
	}

	forbiddenLicensesOnly = dependencies.DependencyInfo{
		Dependencies: []dependencies.Dependency{
			{
				Name:     "library1",
				Version:  "1.0.2",
				Licenses: []string{detectlicense.AGPL1Only.Name},
			},
		},
		Licenses: map[string]string{},
	}

	unrestrictedLicensesOnly = dependencies.DependencyInfo{
		Dependencies: []dependencies.Dependency{
			{
				Name:     "library1",
				Version:  "1.0.2",
				Licenses: []string{detectlicense.MIT.Name},
			},
		},
		Licenses: map[string]string{},
	}

	licensesForAmbassadorServersOnly = dependencies.DependencyInfo{
		Dependencies: []dependencies.Dependency{
			{
				Name:     "library1",
				Version:  "1.0.2",
				Licenses: []string{detectlicense.GPL3Only.Name},
			},
		},
		Licenses: map[string]string{},
	}

	mixOfLicensesIncludingForbidden = dependencies.DependencyInfo{
		Dependencies: []dependencies.Dependency{
			{
				Name:     "library1",
				Version:  "1.0.2",
				Licenses: []string{detectlicense.GPL3Only.Name},
			},
			{
				Name:     "library2",
				Version:  "3.1.3",
				Licenses: []string{detectlicense.MIT.Name},
			},
			{
				Name:     "library3",
				Version:  "1.3.5",
				Licenses: []string{detectlicense.AGPL1Only.Name},
			},
		},
		Licenses: map[string]string{},
	}

	mixOfLicensesWithoutForbidden = dependencies.DependencyInfo{
		Dependencies: []dependencies.Dependency{
			{
				Name:     "library1",
				Version:  "1.0.2",
				Licenses: []string{detectlicense.GPL3Only.Name, detectlicense.MIT.Name},
			},
		},
		Licenses: map[string]string{},
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
			err := testCase.dependencies.UpdateLicenseList()
			require.NoError(t, err)

			require.Equal(t, testCase.expectedLicenses, testCase.dependencies.Licenses)
		})
	}
}

func TestCheckLicensesValidatesAllowedLicenseCorrectly(t *testing.T) {
	testCases := []struct {
		Name            string
		dependencies    dependencies.DependencyInfo
		allowedLicenses detectlicense.LicenseRestriction
	}{
		{
			"Empty dependency list is always allowed",
			emptyDependencies,
			detectlicense.Unrestricted,
		},
		{
			"Unrestricted licenses are OK on Ambassador Labs servers",
			unrestrictedLicensesOnly,
			detectlicense.AmbassadorServers,
		},
		{
			"Unrestricted licenses are OK everywhere",
			unrestrictedLicensesOnly,
			detectlicense.Unrestricted,
		},
		{
			"Restricted licenses are OK on Ambassador Labs servers",
			licensesForAmbassadorServersOnly,
			detectlicense.AmbassadorServers,
		},
		{
			"Mix of licenses without forbidden is allowed on Ambassador Labs servers",
			mixOfLicensesWithoutForbidden,
			detectlicense.AmbassadorServers,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			err := testCase.dependencies.CheckLicenses(testCase.allowedLicenses)
			require.NoError(t, err)
		})
	}
}

func TestCheckLicensesValidatesForbiddenLicensesCorrectly(t *testing.T) {
	testCases := []struct {
		Name            string
		dependencies    dependencies.DependencyInfo
		allowedLicenses detectlicense.LicenseRestriction
	}{
		{
			"It's not possible to allow the use of forbidden licenses by mistake",
			unrestrictedLicensesOnly,
			detectlicense.Forbidden,
		},
		{
			"Forbidden licenses are not allowed on Ambassador Labs servers",
			forbiddenLicensesOnly,
			detectlicense.AmbassadorServers,
		},
		{
			"Forbidden licenses are not allowed on customer servers",
			forbiddenLicensesOnly,
			detectlicense.Unrestricted,
		},
		{
			"Restricted licenses are not OK on customer servers",
			licensesForAmbassadorServersOnly,
			detectlicense.Unrestricted,
		},
		{
			"Mix of licenses including forbidden is rejected in Ambassador Labs servers",
			mixOfLicensesIncludingForbidden,
			detectlicense.Unrestricted,
		},
		{
			"Mix of licenses without forbidden is rejected on customer servers",
			mixOfLicensesWithoutForbidden,
			detectlicense.Unrestricted,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			err := testCase.dependencies.CheckLicenses(testCase.allowedLicenses)
			require.Error(t, err)
		})
	}
}
