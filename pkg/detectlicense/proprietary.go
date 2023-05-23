package detectlicense

import (
	"os"

	"gopkg.in/yaml.v3"
)

//nolint:gochecknoglobals // Would be 'const'.
var ambassadorPrivateRepos = []string{
	"github.com/datawire/telepresence2-proprietary/",
	"github.com/datawire/saas_app/",
	"github.com/datawire/telepresence-pro/",
}

type AmbassadorProprietarySoftware []string

func GetAmbassadorProprietarySoftware(proprietarySoftware ...string) AmbassadorProprietarySoftware {
	ambProprietarySoftware := AmbassadorProprietarySoftware{}
	ambProprietarySoftware = append(ambProprietarySoftware, ambassadorPrivateRepos...)
	ambProprietarySoftware = append(ambProprietarySoftware, proprietarySoftware...)
	return ambProprietarySoftware
}

func (a AmbassadorProprietarySoftware) IsProprietarySoftware(packageName string) bool {
	for _, proprietarySoftware := range a {
		if packageName == proprietarySoftware {
			return true
		}
	}
	return false
}

func (a *AmbassadorProprietarySoftware) ReadProprietarySoftwareFile(name string) error {
	data, err := os.ReadFile(name)
	if err != nil {
		return err
	}

	var proprietary_software []string
	if err = yaml.Unmarshal(data, &proprietary_software); err != nil {
		return err
	}

	*a = append(*a, proprietary_software...)
	return nil
}
