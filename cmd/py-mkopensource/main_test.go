package main

import (
	"encoding/json"
	"github.com/datawire/go-mkopensource/pkg/dependencies"
	"github.com/stretchr/testify/require"
	"io"
	"os"
	"path"
	"testing"
)

func TestMarkdownOutput(t *testing.T) {
	testCases := []struct {
		testName        string
		dependencies    string
		expectedOutput  string
		applicationType ApplicationType
	}{
		{
			"Different dependencies are processed correctly",
			"./testdata/successful-generation/dependency_list.txt",
			"./testdata/successful-generation/expected_markdown.txt",
			internalApplication,
		},
		{
			"Same dependency twice with different version",
			"./testdata/two-versions-of-a-dependency/dependency_list.txt",
			"./testdata/two-versions-of-a-dependency/expected_markdown.txt",
			internalApplication,
		},
		{
			"GPL licenses are allowed for internal use",
			"./testdata/gpl-license/dependency_list.txt",
			"./testdata/gpl-license/expected_markdown.txt",
			internalApplication,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			//Arrange
			pipDependencies, err := os.Open(testCase.dependencies)
			require.NoError(t, err)
			defer func() { _ = pipDependencies.Close() }()

			r, w, pipeErr := os.Pipe()
			require.NoError(t, pipeErr)

			// Act
			err = Main(markdownOutputType, testCase.applicationType, pipDependencies, w)
			require.NoError(t, err)
			_ = w.Close()

			// Assert
			programOutput, readErr := io.ReadAll(r)
			require.NoError(t, readErr)

			expectedOutput := getFileContents(t, testCase.expectedOutput)
			require.Equal(t, expectedOutput, string(programOutput))
		})
	}
}

func TestJsonOutput(t *testing.T) {
	testCases := []struct {
		testName        string
		dependencies    string
		expectedOutput  string
		applicationType ApplicationType
	}{
		{
			"Different dependencies are processed correctly",
			"./testdata/successful-generation/dependency_list.txt",
			"./testdata/successful-generation/expected_json.json",
			internalApplication,
		},
		{
			"Same dependency twice with different version",
			"./testdata/two-versions-of-a-dependency/dependency_list.txt",
			"./testdata/two-versions-of-a-dependency/expected_json.json",
			internalApplication,
		},
		{
			"GPL licenses are allowed for internal use",
			"./testdata/gpl-license/dependency_list.txt",
			"./testdata/gpl-license/expected_json.json",
			internalApplication,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			//Arrange
			pipDependencies, err := os.Open(testCase.dependencies)
			require.NoError(t, err)
			defer func() { _ = pipDependencies.Close() }()

			r, w, pipeErr := os.Pipe()
			require.NoError(t, pipeErr)

			// Act
			err = Main(jsonOutputType, testCase.applicationType, pipDependencies, w)
			require.NoError(t, err)
			_ = w.Close()

			// Assert
			programOutput := getDependencyInfoFromReader(t, r)
			expectedOutput := getDependencyInfoFromFile(t, testCase.expectedOutput)
			require.Equal(t, expectedOutput, programOutput)
		})
	}
}

func TestLicenseErrors(t *testing.T) {
	testCases := []struct {
		testName             string
		dependencies         string
		outputType           OutputType
		applicationType      ApplicationType
		expectedErrorMessage string
	}{
		{
			"GPL licenses are forbidden for external use - Markdown format",
			"./testdata/gpl-license",
			markdownOutputType,
			externalApplication,
			"Dependency 'docutils@0.17.1' uses license 'GNU General Public License v3.0 or later' which is not allowed on applications that run on customer machines.",
		},
		{
			"GPL licenses are forbidden for external use - JSON format",
			"./testdata/gpl-license",
			jsonOutputType,
			externalApplication,
			"Dependency 'docutils@0.17.1' uses license 'GNU General Public License v3.0 or later' which is not allowed on applications that run on customer machines.",
		},
		{
			"AGPL licenses are forbidden for internal use - Markdown format",
			"./testdata/agpl-license",
			markdownOutputType,
			internalApplication,
			"Dependency 'infomap@2.0.2' uses license 'GNU Affero General Public License v3.0 or later' which is forbidden",
		},
		{
			"AGPL licenses are forbidden for internal use - JSON format",
			"./testdata/agpl-license",
			jsonOutputType,
			internalApplication,
			"Dependency 'infomap@2.0.2' uses license 'GNU Affero General Public License v3.0 or later' which is forbidden",
		},
		{
			"AGPL licenses are forbidden for external use - Markdown format",
			"./testdata/agpl-license",
			markdownOutputType,
			externalApplication,
			"is forbidden",
		},
		{
			"AGPL licenses are forbidden for external use - JSON format",
			"./testdata/agpl-license",
			jsonOutputType,
			externalApplication,
			"is forbidden",
		},
		{
			"Unknown licenses are identified correctly",
			"./testdata/unknown-license",
			jsonOutputType,
			externalApplication,
			"Package\"CacheControl\" \"1.99.6\": Could not parse license-string \"UNKNOWN\"",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			//Arrange
			pipDependencies, err := os.Open(path.Join(testCase.dependencies, "dependency_list.txt"))
			require.NoError(t, err)
			defer func() { _ = pipDependencies.Close() }()

			_, w, pipeErr := os.Pipe()
			require.NoError(t, pipeErr)

			// Act
			err = Main(markdownOutputType, testCase.applicationType, pipDependencies, w)
			require.Error(t, err)
			expectedError := getFileContents(t, path.Join(testCase.dependencies, "expected_err.txt"))
			require.Equal(t, expectedError, err.Error())
			_ = w.Close()
		})
	}
}

func getFileContents(t *testing.T, path string) string {
	content, err := os.ReadFile(path)
	if err != nil && err != io.EOF {
		require.NoError(t, err)
	}
	return string(content)
}

func getDependencyInfoFromFile(t *testing.T, path string) *dependencies.DependencyInfo {
	f, err := os.Open(path)
	require.NoError(t, err)

	return getDependencyInfoFromReader(t, f)
}

func getDependencyInfoFromReader(t *testing.T, r io.Reader) *dependencies.DependencyInfo {
	data, readErr := io.ReadAll(r)
	require.NoError(t, readErr)

	jsonOutput := &dependencies.DependencyInfo{}
	err := json.Unmarshal(data, jsonOutput)
	require.NoError(t, err)

	return jsonOutput
}
