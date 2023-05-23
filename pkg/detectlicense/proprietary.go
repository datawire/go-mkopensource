package detectlicense

import (
	"os"

	"gopkg.in/yaml.v3"
)

//nolint:gochecknoglobals // Would be 'const'.
var ambassadorPrivateRepos = map[string]struct{}{
	"github.com/datawire/telepresence2-proprietary/": {},
	"github.com/datawire/saas_app/":                  {},
	"github.com/datawire/telepresence-pro/":          {},
}

type AmbassadorProprietarySoftware map[string]struct{}

func GetAmbassadorProprietarySoftware(proprietarySoftware ...string) AmbassadorProprietarySoftware {
	ambProprietarySoftware := AmbassadorProprietarySoftware{}

	for k, v := range ambassadorPrivateRepos {
		ambProprietarySoftware[k] = v
	}

	for _, v := range proprietarySoftware {
		ambProprietarySoftware[v] = struct{}{}
	}

	return ambProprietarySoftware
}

func (a AmbassadorProprietarySoftware) IsProprietarySoftware(packageName string) bool {
	_, ok := a[packageName]
	return ok
}

func (a AmbassadorProprietarySoftware) ReadProprietarySoftwareFile(name string) error {
	data, err := os.ReadFile(name)
	if err != nil {
		return err
	}

	var proprietarySoftware []string
	if err = yaml.Unmarshal(data, &proprietarySoftware); err != nil {
		return err
	}

	for _, v := range proprietarySoftware {
		a[v] = struct{}{}
	}
	return nil
}
