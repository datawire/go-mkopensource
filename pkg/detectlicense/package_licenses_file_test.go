package detectlicense

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestReadPackageLicensesFromFile(t *testing.T) {
	tmp := t.TempDir()
	pls := map[string][]string{
		"github.com/alpha/one": {"Apache-2.0"},
		"github.com/beta/two":  {"MIT", "LGPL-3.0-only"},
	}
	fn := filepath.Join(tmp, "unparsable-packages.yaml")
	f, err := os.Create(fn)
	require.NoError(t, err)
	require.NoError(t, yaml.NewEncoder(f).Encode(pls))
	require.NoError(t, f.Close())

	plm, err := ReadPackageLicensesFromFile(fn)
	require.NoError(t, err)
	require.Equal(t, map[string]map[License]struct{}{
		"github.com/alpha/one": {
			Apache2: {},
		},
		"github.com/beta/two": {
			MIT:       {},
			LGPL3Only: {},
		},
	}, plm)
}

func TestReadPackageLicensesFromFile_invalidName(t *testing.T) {
	tmp := t.TempDir()
	pls := map[string][]string{
		"github.com/alpha/one": {"invalid-spdx-id"},
	}
	fn := filepath.Join(tmp, "unparsable-packages.yaml")
	f, err := os.Create(fn)
	require.NoError(t, err)
	require.NoError(t, yaml.NewEncoder(f).Encode(pls))
	require.NoError(t, f.Close())

	_, err = ReadPackageLicensesFromFile(fn)
	require.Error(t, err)
}
