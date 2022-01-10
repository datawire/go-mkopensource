package main_test

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	main "github.com/datawire/go-mkopensource"
)

func TestGold(t *testing.T) {
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

			expectedError := getExpectedError(t)
			expectedOutput := getExpectedOutput(t)

			originalStdOut, r, w := interceptStdOut()
			defer func() {
				os.Stdout = originalStdOut
			}()

			actErr := main.Main(&main.CLIArgs{
				OutputFormat:  "txt",
				GoTarFilename: filepath.Join("..", "go1.17.3-testdata.src.tar.gz"),
				Package:       "mod",
			})

			_ = w.Close()

			if expectedError == nil {
				require.NoError(t, actErr)

				programOutput, readErr := io.ReadAll(r)
				require.NoError(t, readErr)
				assert.Equal(t, expectedOutput, programOutput)
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

func interceptStdOut() (originalStdOut *os.File, r *os.File, w *os.File) {
	originalStdOut = os.Stdout

	r, w, _ = os.Pipe()
	os.Stdout = w
	return originalStdOut, r, w
}

func getExpectedError(t *testing.T) []byte {
	expErr, err := os.ReadFile("expected_err.txt")
	if err != nil && !os.IsNotExist(err) {
		require.NoError(t, err)
	}
	return expErr
}

func getExpectedOutput(t *testing.T) []byte {
	expectedOutput, err := os.ReadFile("expected_output.txt")
	if err != nil && !os.IsNotExist(err) {
		require.NoError(t, err)
	}
	return expectedOutput
}
