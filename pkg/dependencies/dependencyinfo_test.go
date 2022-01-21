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
				Licenses: []string{detectlicense.GPL3.Name, detectlicense.BSD2.Name},
			},
		},
		Licenses: map[string]string{},
	}

	dependenciesWithOverlappingLicenses = dependencies.DependencyInfo{
		Dependencies: []dependencies.Dependency{
			{
				Name:     "library1",
				Version:  "1.0.2",
				Licenses: []string{detectlicense.GPL3.Name, detectlicense.BSD2.Name},
			},
			{
				Name:     "library2",
				Version:  "3.1.2",
				Licenses: []string{detectlicense.BSD2.Name},
			},
			{
				Name:     "library2",
				Version:  "3.1.2",
				Licenses: []string{detectlicense.Apache2.Name, detectlicense.GPL3.Name},
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
				detectlicense.GPL3.Name: detectlicense.GPL3.URL,
				detectlicense.BSD2.Name: detectlicense.BSD2.URL,
			},
		},
		{
			"Dependencies with overlapping licenses",
			dependenciesWithOverlappingLicenses,
			map[string]string{
				detectlicense.GPL3.Name:    detectlicense.GPL3.URL,
				detectlicense.BSD2.Name:    detectlicense.BSD2.URL,
				detectlicense.Apache2.Name: detectlicense.Apache2.URL,
				detectlicense.GPL3.Name:    detectlicense.GPL3.URL,
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
			// No licenses does not generate an error
			err := testCase.dependencies.UpdateLicenseList()
			require.NoError(t, err)

			require.Equal(t, testCase.expectedLicenses, testCase.dependencies.Licenses)
		})
	}
}
