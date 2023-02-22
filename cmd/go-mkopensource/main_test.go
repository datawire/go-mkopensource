package main_test

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	main "github.com/datawire/go-mkopensource/cmd/go-mkopensource"
	"github.com/datawire/go-mkopensource/pkg/dependencies"
)

func TestSuccessfulMarkdownOutput(t *testing.T) {
	testCases := []struct {
		testName                string
		testData                string
		applicationType         string
		supportedGoVersionRegEx string
	}{
		{
			testName:                "01-intern-new - markdown output",
			testData:                "testdata/01-intern-new",
			applicationType:         "external",
			supportedGoVersionRegEx: `.*`,
		},
		{
			testName:                "02-replace - markdown output",
			testData:                "testdata/02-replace",
			applicationType:         "external",
			supportedGoVersionRegEx: `.*`,
		},
		{
			testName:                "04-nodeps - markdown output",
			testData:                "testdata/04-nodeps",
			applicationType:         "external",
			supportedGoVersionRegEx: `.*`,
		},
		{
			testName:                "05-subpatent - markdown output",
			testData:                "testdata/05-subpatent",
			applicationType:         "external",
			supportedGoVersionRegEx: `.*`,
		},
		{
			testName:                "One dependency with multiple licenses",
			testData:                "testdata/06-multiple-licenses",
			applicationType:         "external",
			supportedGoVersionRegEx: `.*`,
		},
		{
			testName:                "GPL license is allowed for internal use",
			testData:                "testdata/08-allowed-for-internal-use-only",
			applicationType:         "internal",
			supportedGoVersionRegEx: `.*`,
		},
		{
			testName:                "09-out-of-date-dependencies - Dependency not found",
			testData:                "testdata/09-out-of-date-dependencies-markdown",
			applicationType:         "external",
			supportedGoVersionRegEx: `^go1\.1[89]\..*`,
		},
		{
			testName:                "11-dependency-missing-from-go-mod - Dependency missing from go.mod is added",
			testData:                "testdata/11-dependency-missing-from-go-mod",
			applicationType:         "external",
			supportedGoVersionRegEx: `.*`,
		},
	}

	workingDir := getWorkingDir(t)

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			re := regexp.MustCompile(testCase.supportedGoVersionRegEx)
			if !re.Match([]byte(runtime.Version())) {
				t.Skipf("Test does not support go version %s", runtime.Version())
			}

			defer func() {
				require.NoError(t, os.Chdir(workingDir))
			}()

			require.NoError(t, os.Chdir(testCase.testData))

			originalStdOut, r, w := interceptStdOut()
			defer func() {
				os.Stdout = originalStdOut
			}()

			actErr := main.Main(&main.CLIArgs{
				OutputFormat:    "txt",
				GoTarFilename:   filepath.Join("..", "go1.17.3-testdata.src.tar.gz"),
				Package:         "mod",
				OutputType:      "markdown",
				ApplicationType: testCase.applicationType,
			})

			_ = w.Close()

			require.NoError(t, actErr)

			programOutput, readErr := io.ReadAll(r)
			require.NoError(t, readErr)

			expectedOutput := getFileContents(t, "expected_markdown_output.txt")

			assert.Equal(t, string(expectedOutput), string(programOutput))
		})
	}
}

func TestSuccessfulJsonOutput(t *testing.T) {
	testCases := []struct {
		testName                string
		testData                string
		applicationType         string
		supportedGoVersionRegEx string
	}{
		{
			testName:                "01-intern-new",
			testData:                "testdata/01-intern-new",
			applicationType:         "external",
			supportedGoVersionRegEx: `.*`,
		},
		{
			testName:                "02-replace",
			testData:                "testdata/02-replace",
			applicationType:         "external",
			supportedGoVersionRegEx: `.*`,
		},
		{
			testName:                "04-nodeps",
			testData:                "testdata/04-nodeps",
			applicationType:         "external",
			supportedGoVersionRegEx: `.*`,
		},
		{
			testName:                "05-subpatent",
			testData:                "testdata/05-subpatent",
			applicationType:         "external",
			supportedGoVersionRegEx: `.*`,
		},
		{
			testName:                "One dependency with multiple licenses",
			testData:                "testdata/06-multiple-licenses",
			applicationType:         "external",
			supportedGoVersionRegEx: `.*`,
		},
		{
			testName:                "GPL license is allowed for internal use",
			testData:                "testdata/08-allowed-for-internal-use-only",
			applicationType:         "internal",
			supportedGoVersionRegEx: `.*`,
		},
		{
			testName:                "09-out-of-date-dependencies - Dependency not found",
			testData:                "testdata/09-out-of-date-dependencies-json",
			applicationType:         "external",
			supportedGoVersionRegEx: `^go1\.1[89]\..*`,
		},
	}

	workingDir := getWorkingDir(t)

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			re := regexp.MustCompile(testCase.supportedGoVersionRegEx)
			if !re.Match([]byte(runtime.Version())) {
				t.Skipf("Test does not support go version %s", runtime.Version())
			}

			defer func() {
				require.NoError(t, os.Chdir(workingDir))
			}()

			require.NoError(t, os.Chdir(testCase.testData))

			originalStdOut, r, w := interceptStdOut()
			defer func() {
				os.Stdout = originalStdOut
			}()

			actErr := main.Main(&main.CLIArgs{
				OutputFormat:    "txt",
				GoTarFilename:   filepath.Join("..", "go1.17.3-testdata.src.tar.gz"),
				Package:         "mod",
				OutputType:      "json",
				ApplicationType: testCase.applicationType,
			})

			_ = w.Close()

			require.NoError(t, actErr)

			jsonOutput := getDependencyInfoFromReader(t, r)
			expectedJson := getDependencyInfoFromFile(t, "expected_json_output.json")

			assert.Equal(t, expectedJson, jsonOutput)
		})
	}
}

func TestSuccessfulTarOutput(t *testing.T) {
	testCases := []struct {
		testName                string
		testData                string
		applicationType         string
		supportedGoVersionRegEx string
	}{
		{
			testName:                "01-intern-new - markdown output",
			testData:                "testdata/01-intern-new",
			applicationType:         "external",
			supportedGoVersionRegEx: `.*`,
		},
		{
			testName:                "02-replace - markdown output",
			testData:                "testdata/02-replace",
			applicationType:         "external",
			supportedGoVersionRegEx: `.*`,
		},
		{
			testName:                "04-nodeps - markdown output",
			testData:                "testdata/04-nodeps",
			applicationType:         "external",
			supportedGoVersionRegEx: `.*`,
		},
		{
			testName:                "05-subpatent - markdown output",
			testData:                "testdata/05-subpatent",
			applicationType:         "external",
			supportedGoVersionRegEx: `.*`,
		},
		{
			testName:                "One dependency with multiple licenses",
			testData:                "testdata/06-multiple-licenses",
			applicationType:         "external",
			supportedGoVersionRegEx: `.*`,
		},
		{
			testName:                "GPL license is allowed for internal use",
			testData:                "testdata/08-allowed-for-internal-use-only",
			applicationType:         "internal",
			supportedGoVersionRegEx: `.*`,
		},
		{
			testName:                "09-out-of-date-dependencies - Dependency not found",
			testData:                "testdata/09-out-of-date-dependencies-tar",
			applicationType:         "external",
			supportedGoVersionRegEx: `^go1\.1[89]\..*`,
		},
	}

	workingDir := getWorkingDir(t)

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			re := regexp.MustCompile(testCase.supportedGoVersionRegEx)
			if !re.Match([]byte(runtime.Version())) {
				t.Skipf("Test does not support go version %s", runtime.Version())
			}

			defer func() {
				require.NoError(t, os.Chdir(workingDir))
			}()

			require.NoError(t, os.Chdir(testCase.testData))

			originalStdOut, r, w := interceptStdOut()
			defer func() {
				os.Stdout = originalStdOut
			}()

			actErr := main.Main(&main.CLIArgs{
				OutputFormat:    "tar",
				OutputName:      "",
				GoTarFilename:   filepath.Join("..", "go1.17.3-testdata.src.tar.gz"),
				Package:         "mod",
				ApplicationType: testCase.applicationType,
			})

			_ = w.Close()

			require.NoError(t, actErr)

			fileContents, err := listTarContents(t, r)
			require.NoError(t, err)

			expectedTarContents := getFileContents(t, "expected_tar_contents.txt")

			assert.Equal(t, string(expectedTarContents), fileContents)
		})
	}
}

func TestErrorScenarios(t *testing.T) {
	testCases := []struct {
		testName                string
		testData                string
		packagesFlag            string
		outputTypeFlag          string
		supportedGoVersionRegEx string
	}{
		{
			testName:                "testdata/00-intern-old",
			testData:                "testdata/00-intern-old",
			packagesFlag:            "mod",
			outputTypeFlag:          "full",
			supportedGoVersionRegEx: `.*`,
		},
		{
			testName:                "Multiple errors",
			testData:                "testdata/03-multierror",
			packagesFlag:            "mod",
			outputTypeFlag:          "full",
			supportedGoVersionRegEx: `.*`,
		},
		{
			testName:                "Forbidden license",
			testData:                "testdata/07-forbidden-license",
			packagesFlag:            "mod",
			outputTypeFlag:          "full",
			supportedGoVersionRegEx: `.*`,
		},
		{
			testName:                "License not allowed on distributed applications",
			testData:                "testdata/08-allowed-for-internal-use-only",
			packagesFlag:            "mod",
			outputTypeFlag:          "full",
			supportedGoVersionRegEx: `.*`,
		},
		{
			testName:                "Can't update dependencies due to Go version 1.16",
			testData:                "testdata/09-out-of-date-dependencies-markdown",
			packagesFlag:            "mod",
			outputTypeFlag:          "external",
			supportedGoVersionRegEx: `^go1\.16\..*`,
		},
		{
			testName:                "Can't update dependencies due to Go version 1.17",
			testData:                "testdata/09-out-of-date-dependencies-json",
			packagesFlag:            "mod",
			outputTypeFlag:          "external",
			supportedGoVersionRegEx: `^go1\.17\..*`,
		},
		{
			testName:                "Removed dependency is required by the application",
			testData:                "testdata/10-removed-dependency-imported-by-app",
			packagesFlag:            "mod",
			outputTypeFlag:          "external",
			supportedGoVersionRegEx: `.*`,
		},
	}

	workingDir := getWorkingDir(t)

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			re := regexp.MustCompile(testCase.supportedGoVersionRegEx)
			if !re.Match([]byte(runtime.Version())) {
				t.Skipf("Test does not support go version %s", runtime.Version())
			}

			defer func() {
				require.NoError(t, os.Chdir(workingDir))
			}()

			require.NoError(t, os.Chdir(testCase.testData))

			expectedError := getFileContents(t, "expected_err.txt")

			actErr := main.Main(&main.CLIArgs{
				OutputFormat:    "txt",
				GoTarFilename:   filepath.Join("..", "go1.17.3-testdata.src.tar.gz"),
				Package:         testCase.packagesFlag,
				OutputType:      testCase.outputTypeFlag,
				ApplicationType: "external",
			})

			if assert.Error(t, actErr) {
				// Use this instead of assert.EqualError so that we diff
				// output, which is helpful for long strings.
				assert.Equal(t, strings.TrimSpace(string(expectedError)), strings.TrimSpace(actErr.Error()))
			}
		})
	}
}

func getWorkingDir(t *testing.T) string {
	workingDir, err := os.Getwd()
	require.NoError(t, err)
	return workingDir
}

func listTarContents(t *testing.T, r *os.File) (string, error) {
	gzipFile, err := gzip.NewReader(r)
	if err != nil {
		return "", err
	}

	tarFile := tar.NewReader(gzipFile)
	files := []string{}
	for {
		header, err := tarFile.Next()
		if err == io.EOF {
			break // End of archive
		}

		require.NoError(t, err)
		files = append(files, header.Name)
	}

	fileContents := strings.Join(files, "\n")
	return fileContents, nil
}

func interceptStdOut() (originalStdOut *os.File, r *os.File, w *os.File) {
	originalStdOut = os.Stdout

	r, w, _ = os.Pipe()
	os.Stdout = w
	return originalStdOut, r, w
}

func getFileContents(t *testing.T, path string) []byte {
	expErr, err := os.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		require.NoError(t, err)
	}
	return expErr
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
	err := jsonOutput.Unmarshal(data)
	require.NoError(t, err)

	return jsonOutput
}
