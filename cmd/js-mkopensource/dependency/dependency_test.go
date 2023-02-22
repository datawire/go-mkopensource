package dependency_test

import (
	"io"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/datawire/go-mkopensource/cmd/js-mkopensource/dependency"
	"github.com/datawire/go-mkopensource/pkg/dependencies"
	"github.com/datawire/go-mkopensource/pkg/detectlicense"
)

func TestSuccessfulGeneration(t *testing.T) {
	testCases := []struct {
		testName string
		input    string
	}{
		{
			"Dependency identifier in the format @name@version",
			"./testdata/dependency-with-special-characters",
		},
		{
			"Multiple dependencies",
			"./testdata/multiple-licenses",
		},
		{
			"One dependency with multiple licenses",
			"./testdata/dependencies-with-two-licenses",
		},
		{
			"Hardcoded dependencies are properly parsed",
			"./testdata/hardcoded-dependencies",
		},
		{
			"GPL license is allowed in internal software",
			"./testdata/dependency-with-gpl-license",
		},
		{
			"License field can be string or array",
			"./testdata/license-field-as-an-array",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			//Arrange
			nodeDependencies := getNodeDependencies(t, path.Join(testCase.input, "dependencies.json"))
			defer func() { _ = nodeDependencies.Close() }()

			// Act
			dependencyInformation, err := dependency.GetDependencyInformation(nodeDependencies, detectlicense.AmbassadorServers)
			require.NoError(t, err)

			// Assert
			expectedJson := getDependencyInfoFromFile(t, path.Join(testCase.input, "expected_output.json"))
			require.Equal(t, *expectedJson, dependencyInformation)
		})
	}
}

func TestErrorScenarios(t *testing.T) {
	testCases := []struct {
		testName string
		input    string
	}{
		{
			"Invalid Json input",
			"./testdata/invalid-json",
		},
		{
			"Unknown license identifier",
			"./testdata/unknown-license",
		},
		{
			"Missing license",
			"./testdata/missing-license",
		},
		{
			"Empty license field",
			"./testdata/empty-license",
		},
		{
			"Hardcode dependency with different version is rejected",
			"./testdata/hardcoded-dependencies-but-different-version",
		},
		{
			"GPL license is not allowed in distributed software",
			"./testdata/dependency-with-gpl-license",
		},
		{
			"AGPL license is forbidden",
			"./testdata/dependency-with-agpl-license",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			//Arrange
			nodeDependencies := getNodeDependencies(t, path.Join(testCase.input, "dependencies.json"))
			defer func() { _ = nodeDependencies.Close() }()

			// Act
			_, err := dependency.GetDependencyInformation(nodeDependencies, detectlicense.Unrestricted)

			// Assert
			require.Error(t, err)
			expectedError := getFileContents(t, path.Join(testCase.input, "expected_err.txt"))
			assert.Equal(t, string(expectedError), err.Error())
		})
	}
}

func getNodeDependencies(t *testing.T, dependencyFile string) *os.File {
	nodeDependencies, openErr := os.Open(dependencyFile)
	require.NoError(t, openErr)
	return nodeDependencies
}

func getDependencyInfoFromFile(t *testing.T, path string) *dependencies.DependencyInfo {
	f, openErr := os.Open(path)
	require.NoError(t, openErr)

	data, readErr := io.ReadAll(f)
	require.NoError(t, readErr)

	dependencyInfo := &dependencies.DependencyInfo{}
	unmarshalErr := dependencyInfo.Unmarshal(data)
	require.NoError(t, unmarshalErr)

	return dependencyInfo
}

func getFileContents(t *testing.T, path string) []byte {
	expErr, err := os.ReadFile(path)
	if err != nil {
		require.NoError(t, err)
	}
	return expErr
}
