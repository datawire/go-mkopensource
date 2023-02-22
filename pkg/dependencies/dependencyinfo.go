package dependencies

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/datawire/go-mkopensource/pkg/detectlicense"
	"github.com/datawire/go-mkopensource/pkg/util"
)

//nolint:gochecknoglobals // Would be 'const'.
var licensesByName = func() map[string]detectlicense.License {
	var licenses = []detectlicense.License{
		detectlicense.AmbassadorProprietary,
		detectlicense.ZeroBSD,
		detectlicense.Apache2,
		detectlicense.AFL21,
		detectlicense.AGPL1Only,
		detectlicense.AGPL1OrLater,
		detectlicense.AGPL3Only,
		detectlicense.AGPL3OrLater,
		detectlicense.BSD1,
		detectlicense.BSD2,
		detectlicense.BSD3,
		detectlicense.CcBy30,
		detectlicense.CcBy40,
		detectlicense.CcBySa40,
		detectlicense.Cc010,
		detectlicense.EPL10,
		detectlicense.GPL1Only,
		detectlicense.GPL1OrLater,
		detectlicense.GPL2Only,
		detectlicense.GPL2OrLater,
		detectlicense.GPL3Only,
		detectlicense.GPL3OrLater,
		detectlicense.ISC,
		detectlicense.LGPL2Only,
		detectlicense.LGPL2OrLater,
		detectlicense.LGPL21Only,
		detectlicense.LGPL21OrLater,
		detectlicense.LGPL3Only,
		detectlicense.LGPL3OrLater,
		detectlicense.MIT,
		detectlicense.MPL11,
		detectlicense.MPL2,
		detectlicense.ODCBy10,
		detectlicense.OFL11,
		detectlicense.Python20,
		detectlicense.PSF,
		detectlicense.PublicDomain,
		detectlicense.Unicode2015,
		detectlicense.Unlicense,
		detectlicense.WTFPL,
	}
	var ret = make(map[string]detectlicense.License, len(licenses))
	for _, lic := range licenses {
		ret[lic.Name] = lic
	}
	return ret
}()

type Dependency struct {
	Name     string
	Version  string
	Licenses util.Set[detectlicense.License]
}

var (
	_ json.Marshaler   = Dependency{}
	_ json.Unmarshaler = (*Dependency)(nil)
)

func (d Dependency) MarshalJSON() ([]byte, error) {
	raw := struct {
		Name     string   `json:"name"`
		Version  string   `json:"version"`
		Licenses []string `json:"licenses"`
	}{
		Name:    d.Name,
		Version: d.Version,
	}
	for lic := range d.Licenses {
		raw.Licenses = append(raw.Licenses, lic.Name)
	}
	sort.Strings(raw.Licenses)
	return json.Marshal(raw)
}

func (d *Dependency) UnmarshalJSON(dat []byte) error {
	var raw struct {
		Name     string   `json:"name"`
		Version  string   `json:"version"`
		Licenses []string `json:"licenses"`
	}
	if err := json.Unmarshal(dat, &raw); err != nil {
		return err
	}
	d.Name = raw.Name
	d.Version = raw.Version
	d.Licenses = make(util.Set[detectlicense.License], len(raw.Licenses))
	for _, licName := range raw.Licenses {
		d.Licenses.Insert(licensesByName[licName])
	}
	return nil
}

type DependencyInfo struct {
	Dependencies []Dependency `json:"dependencies"`
}

var (
	_ json.Marshaler = DependencyInfo{}
)

func (di DependencyInfo) MarshalJSON() ([]byte, error) {
	raw := struct {
		Dependencies []Dependency      `json:"dependencies"`
		Licenses     map[string]string `json:"licenseInfo"`
	}{
		Dependencies: di.Dependencies,
		Licenses:     make(map[string]string),
	}
	for _, dep := range di.Dependencies {
		for lic := range dep.Licenses {
			raw.Licenses[lic.Name] = lic.URL
		}
	}
	return json.Marshal(raw)
}

func (dependency Dependency) CheckLicenseRestrictions(licenseRestriction detectlicense.LicenseRestriction) []error {
	var errs []error
	for license := range dependency.Licenses {
		switch {
		case license.Restriction == detectlicense.Forbidden:
			errs = append(errs, fmt.Errorf("Dependency '%s@%s' uses license '%s' which is forbidden.", dependency.Name,
				dependency.Version, license.Name))
		case license.Restriction < licenseRestriction:
			errs = append(errs, fmt.Errorf("Dependency '%s@%s' uses license '%s' which is not allowed on applications that run on customer machines.",
				dependency.Name, dependency.Version, license.Name))
		}
	}
	return errs
}
