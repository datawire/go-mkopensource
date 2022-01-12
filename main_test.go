package main_test

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	main "github.com/datawire/go-mkopensource"
)

func TestGolden(t *testing.T) {
	testCases := []struct {
		testName       string
		testData       string
		packagesFlag   string
		outputTypeFlag string
	}{
		{
			testName:       "",
			testData:       "testdata/01-intern-new",
			packagesFlag:   "mod",
			outputTypeFlag: "markdown",
		},
		{
			testName:       "",
			testData:       "testdata/02-replace",
			packagesFlag:   "mod",
			outputTypeFlag: "markdown",
		},
		{
			testName:       "",
			testData:       "testdata/04-nodeps",
			packagesFlag:   "mod",
			outputTypeFlag: "markdown",
		},
		{
			testName:       "",
			testData:       "testdata/05-subpatent",
			packagesFlag:   "mod",
			outputTypeFlag: "markdown",
		},
		{
			testName:       "",
			testData:       "testdata/06-multiple-licenses",
			packagesFlag:   "mod",
			outputTypeFlag: "markdown",
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
				OutputFormat:  "txt",
				GoTarFilename: filepath.Join("..", "go1.17.3-testdata.src.tar.gz"),
				Package:       testCase.packagesFlag,
				OutputType:    testCase.outputTypeFlag,
			})

			_ = w.Close()

			require.NoError(t, actErr)

			programOutput, readErr := io.ReadAll(r)
			require.NoError(t, readErr)

			expectedOutput := getFileContents(t, "expected_full_output.txt")
			assert.Equal(t, string(expectedOutput), string(programOutput))
		})
	}
}

func TestErrorScenarios(t *testing.T) {
	testCases := []struct {
		testData       string
		packagesFlag   string
		outputTypeFlag string
	}{
		{
			testData:       "testdata/00-intern-old",
			packagesFlag:   "mod",
			outputTypeFlag: "full",
		},
		{
			testData:       "testdata/03-multierror",
			packagesFlag:   "mod",
			outputTypeFlag: "full",
		},
	}

	workingDir := getWorkingDir(t)

	for _, testCase := range testCases {
		t.Run(testCase.testData, func(t *testing.T) {
			defer func() {
				require.NoError(t, os.Chdir(workingDir))
			}()

			require.NoError(t, os.Chdir(testCase.testData))

			expectedError := getFileContents(t, "expected_err.txt")

			actErr := main.Main(&main.CLIArgs{
				OutputFormat:  "txt",
				GoTarFilename: filepath.Join("..", "go1.17.3-testdata.src.tar.gz"),
				Package:       testCase.packagesFlag,
				OutputType:    testCase.outputTypeFlag,
			})

			if assert.Error(t, actErr) {
				// Use this instead of assert.EqualError so that we diff
				// output, which is helpful for long strings.
				assert.Equal(t, strings.TrimSpace(string(expectedError)), actErr.Error())
			}
		})
	}
}

func getWorkingDir(t *testing.T) string {
	workingDir, err := os.Getwd()
	require.NoError(t, err)
	return workingDir
}

func TestLicenseTxtOutput(t *testing.T) {
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
				OutputFormat:  "txt",
				OutputType:    "license",
				GoTarFilename: filepath.Join("..", "go1.17.3-testdata.src.tar.gz"),
				Package:       "mod",
			})

			_ = w.Close()

			if expectedError == nil {
				require.NoError(t, actErr)

				programOutput, readErr := io.ReadAll(r)
				require.NoError(t, readErr)

				expectedOutput := getFileContents(t, "expected_license_output.txt")
				assert.Equal(t, string(expectedOutput), string(programOutput))
			} else {
				if assert.Error(t, actErr) {
					// Use this instead of assert.EqualError so that we diff
					// output, which is helpful for long strings.
					assert.Equal(t, strings.TrimSpace(string(expectedError)), actErr.Error())
				}
			}
		})
	}
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
					assert.Equal(t, strings.TrimSpace(string(expectedError)), actErr.Error())
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
