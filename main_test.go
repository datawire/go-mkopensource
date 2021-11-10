package main_test

import (
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

			expErr, err := os.ReadFile("expected_err.txt")
			if err != nil && !os.IsNotExist(err) {
				require.NoError(t, err)
			}

			actErr := main.Main(&main.CLIArgs{
				OutputFormat:  "txt",
				GoTarFilename: filepath.Join("..", "go1.17.3-testdata.src.tar.gz"),
				Package:       "mod",
			})
			if expErr == nil {
				assert.NoError(t, actErr)
			} else {
				assert.EqualError(t, actErr, strings.TrimSpace(string(expErr)))
			}
		})
	}
}
