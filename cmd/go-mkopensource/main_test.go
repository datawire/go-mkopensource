package main_test

import (
	"archive/tar"
	"compress/gzip"
	"github.com/datawire/go-mkopensource/pkg/dependencies"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	main "github.com/datawire/go-mkopensource/cmd/go-mkopensource"
)

func TestSuccessfulMarkdownOutput(t *testing.T) {
	testCases := []struct {
		testName        string
		testData        string
		applicationType string
	}{
		{
			testName:        "01-intern-new - markdown output",
			testData:        "testdata/01-intern-new",
			applicationType: "external",
		},
		{
			testName:        "02-replace - markdown output",
			testData:        "testdata/02-replace",
			applicationType: "external",
		},
		{
			testName:        "04-nodeps - markdown output",
			testData:        "testdata/04-nodeps",
			applicationType: "external",
		},
		{
			testName:        "05-subpatent - markdown output",
			testData:        "testdata/05-subpatent",
			applicationType: "external",
		},
		{
			testName:        "One dependency with multiple licenses",
			testData:        "testdata/06-multiple-licenses",
			applicationType: "external",
		},
		{
			testName:        "GPL license is allowed for internal use",
			testData:        "testdata/08-allowed-for-internal-use-only",
			applicationType: "internal",
		},
		{
			testName:        "09-out-of-date-dependencies - Dependency not found",
			testData:        "testdata/09-out-of-date-dependencies-markdown",
			applicationType: "external",
		},
	}

	workingDir := getWorkingDir(t)

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
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
		testName        string
		testData        string
		applicationType string
	}{
		{
			testName:        "01-intern-new",
			testData:        "testdata/01-intern-new",
			applicationType: "external",
		},
		{
			testName:        "02-replace",
			testData:        "testdata/02-replace",
			applicationType: "external",
		},
		{
			testName:        "04-nodeps",
			testData:        "testdata/04-nodeps",
			applicationType: "external",
		},
		{
			testName:        "05-subpatent",
			testData:        "testdata/05-subpatent",
			applicationType: "external",
		},
		{
			testName:        "One dependency with multiple licenses",
			testData:        "testdata/06-multiple-licenses",
			applicationType: "external",
		},
		{
			testName:        "GPL license is allowed for internal use",
			testData:        "testdata/08-allowed-for-internal-use-only",
			applicationType: "internal",
		},
		{
			testName:        "09-out-of-date-dependencies - Dependency not found",
			testData:        "testdata/09-out-of-date-dependencies-json",
			applicationType: "external",
		},
	}

	workingDir := getWorkingDir(t)

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
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

func TestErrorScenarios(t *testing.T) {
	testCases := []struct {
		testName       string
		testData       string
		packagesFlag   string
		outputTypeFlag string
	}{
		{
			testName:       "testdata/00-intern-old",
			testData:       "testdata/00-intern-old",
			packagesFlag:   "mod",
			outputTypeFlag: "full",
		},
		{
			testName:       "Multiple errors",
			testData:       "testdata/03-multierror",
			packagesFlag:   "mod",
			outputTypeFlag: "full",
		},
		{
			testName:       "Forbidden license",
			testData:       "testdata/07-forbidden-license",
			packagesFlag:   "mod",
			outputTypeFlag: "full",
		},
		{
			testName:       "License not allowed on distributed applications",
			testData:       "testdata/08-allowed-for-internal-use-only",
			packagesFlag:   "mod",
			outputTypeFlag: "full",
		},
	}

	workingDir := getWorkingDir(t)

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
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

func TestTarOutput(t *testing.T) {
	root, err := os.Getwd()
	require.NoError(t, err)

	direntries, err := os.ReadDir("testdata")
	require.NoError(t, err)
	for _, direntry := range direntries {
		if !direntry.IsDir() {
			continue
		}
		name := direntry.Name()
		t.Run(name, func(t *testing.T) {
			defer func() {
				require.NoError(t, os.Chdir(root))
			}()

			require.NoError(t, os.Chdir(filepath.Join("testdata", name)))

			expectedError := getFileContents(t, "expected_err.txt")

			originalStdOut, r, w := interceptStdOut()
			defer func() {
				os.Stdout = originalStdOut
			}()

			actErr := main.Main(&main.CLIArgs{
				OutputFormat:  "tar",
				OutputName:    "",
				GoTarFilename: filepath.Join("..", "go1.17.3-testdata.src.tar.gz"),
				Package:       "mod",
			})

			_ = w.Close()

			if expectedError == nil {
				require.NoError(t, actErr)

				fileContents, err := listTarContents(t, r)
				require.NoError(t, err)

				expectedTarContents := getFileContents(t, "expected_tar_contents.txt")
				assert.Equal(t, string(expectedTarContents), fileContents)
			} else {
				if assert.Error(t, actErr) {
					// Use this instead of assert.EqualError so that we diff
					// output, which is helpful for long strings.
					assert.Equal(t, strings.TrimSpace(string(expectedError)), strings.TrimSpace(actErr.Error()))
				}
			}
		})
	}
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
