package main

import (
	"encoding/json"
	"github.com/datawire/go-mkopensource/pkg/dependencies"
	"github.com/stretchr/testify/require"
	"io"
	"os"
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
			require.Equal(t, string(expectedOutput), string(programOutput))
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

func TestForbiddenLicenses(t *testing.T) {
	testCases := []struct {
		testName             string
		dependencies         string
		outputType           OutputType
		applicationType      ApplicationType
		expectedErrorMessage string
	}{
		{
			"GPL licenses are forbidden for external use - Markdown format",
			"./testdata/gpl-license/dependency_list.txt",
			markdownOutputType,
			externalApplication,
			"should not be used since it should not run on customer servers",
		},
		{
			"GPL licenses are forbidden for external use - JSON format",
			"./testdata/gpl-license/dependency_list.txt",
			jsonOutputType,
			externalApplication,
			"should not be used since it should not run on customer servers",
		},
		{
			"AGPL licenses are forbidden for internal use - Markdown format",
			"./testdata/agpl-license/dependency_list.txt",
			markdownOutputType,
			internalApplication,
			"is forbidden",
		},
		{
			"AGPL licenses are forbidden for internal use - JSON format",
			"./testdata/agpl-license/dependency_list.txt",
			jsonOutputType,
			internalApplication,
			"is forbidden",
		},
		{
			"AGPL licenses are forbidden for external use - Markdown format",
			"./testdata/agpl-license/dependency_list.txt",
			markdownOutputType,
			externalApplication,
			"is forbidden",
		},
		{
			"AGPL licenses are forbidden for external use - JSON format",
			"./testdata/agpl-license/dependency_list.txt",
			jsonOutputType,
			externalApplication,
			"is forbidden",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			//Arrange
			pipDependencies, err := os.Open(testCase.dependencies)
			require.NoError(t, err)
			defer func() { _ = pipDependencies.Close() }()

			_, w, pipeErr := os.Pipe()
			require.NoError(t, pipeErr)

			// Act
			err = Main(markdownOutputType, testCase.applicationType, pipDependencies, w)
			require.Error(t, err)
			require.Contains(t, err.Error(), testCase.expectedErrorMessage)
			_ = w.Close()
		})
	}
}

func getFileContents(t *testing.T, path string) []byte {
	content, err := os.ReadFile(path)
	if err != nil && err != io.EOF {
		require.NoError(t, err)
	}
	return content
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
