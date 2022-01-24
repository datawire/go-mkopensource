package dependencies

import (
	"encoding/json"
	"fmt"
	"github.com/datawire/go-mkopensource/pkg/detectlicense"
)

//nolint:gochecknoglobals // Can't be a constant
var licensesByName = map[string]detectlicense.License{
	detectlicense.AmbassadorProprietary.Name: detectlicense.AmbassadorProprietary,
	detectlicense.Apache2.Name:               detectlicense.Apache2,
	detectlicense.AGPL1Only.Name:             detectlicense.AGPL1Only,
	detectlicense.AGPL1OrLater.Name:          detectlicense.AGPL1OrLater,
	detectlicense.AGPL3Only.Name:             detectlicense.AGPL3Only,
	detectlicense.AGPL3OrLater.Name:          detectlicense.AGPL3OrLater,
	detectlicense.BSD1.Name:                  detectlicense.BSD1,
	detectlicense.BSD2.Name:                  detectlicense.BSD2,
	detectlicense.BSD3.Name:                  detectlicense.BSD3,
	detectlicense.CcBySa40.Name:              detectlicense.CcBySa40,
	detectlicense.GPL1Only.Name:              detectlicense.GPL1Only,
	detectlicense.GPL1OrLater.Name:           detectlicense.GPL1OrLater,
	detectlicense.GPL2Only.Name:              detectlicense.GPL2Only,
	detectlicense.GPL2OrLater.Name:           detectlicense.GPL2OrLater,
	detectlicense.GPL3Only.Name:              detectlicense.GPL3Only,
	detectlicense.GPL3OrLater.Name:           detectlicense.GPL3OrLater,
	detectlicense.ISC.Name:                   detectlicense.ISC,
	detectlicense.LGPL2Only.Name:             detectlicense.LGPL2Only,
	detectlicense.LGPL2OrLater.Name:          detectlicense.LGPL2OrLater,
	detectlicense.LGPL21Only.Name:            detectlicense.LGPL21Only,
	detectlicense.LGPL21OrLater.Name:         detectlicense.LGPL21OrLater,
	detectlicense.LGPL3Only.Name:             detectlicense.LGPL3Only,
	detectlicense.LGPL3OrLater.Name:          detectlicense.LGPL3OrLater,
	detectlicense.MIT.Name:                   detectlicense.MIT,
	detectlicense.MPL2.Name:                  detectlicense.MPL2,
	detectlicense.PSF.Name:                   detectlicense.PSF,
	detectlicense.PublicDomain.Name:          detectlicense.PublicDomain,
	detectlicense.Unicode2015.Name:           detectlicense.Unicode2015}

type DependencyInfo struct {
	Dependencies []Dependency      `json:"dependencies"`
	Licenses     map[string]string `json:"licenseInfo"`
}

type Dependency struct {
	Name     string   `json:"name"`
	Version  string   `json:"version"`
	Licenses []string `json:"licenses"`
}

func NewDependencyInfo() DependencyInfo {
	return DependencyInfo{
		Dependencies: []Dependency{},
		Licenses:     map[string]string{},
	}
}

func (d *DependencyInfo) Unmarshal(data []byte) error {
	if err := json.Unmarshal(data, d); err != nil {
		return err
	}

	return nil
}

func (d *DependencyInfo) UpdateLicenseList() error {
	usedLicenses := map[string]detectlicense.License{}

	for _, dependency := range d.Dependencies {
		for _, licenseName := range dependency.Licenses {
			license, err := getLicenseFromName(licenseName)
			if err != nil {
				return err
			}
			usedLicenses[license.Name] = license
		}
	}

	for k, v := range usedLicenses {
		d.Licenses[k] = v.URL
	}

	return nil
}

func getLicenseFromName(licenseName string) (detectlicense.License, error) {
	license, ok := licensesByName[licenseName]
	if !ok {
		return detectlicense.License{}, fmt.Errorf("license details for '%s' are not known", licenseName)
	}
	return license, nil
}

// CheckLicenses checks that the licenses used by the dependencies are known and allowed to be used
//in an application based on the buiness logic described here: https://www.notion.so/datawire/License-Management-5194ca50c9684ff4b301143806c92157.
//This function must be called after parsing of the licenses has been done.
func (d *DependencyInfo) CheckLicenses(allowedLicenses detectlicense.AllowedLicenseUse) error {
	if allowedLicenses == detectlicense.Forbidden {
		return fmt.Errorf("forbidden licenses should not be used")
	}

	for _, dependency := range d.Dependencies {
		for _, licenseName := range dependency.Licenses {
			license, err := getLicenseFromName(licenseName)
			if err != nil {
				return err
			}

			if license.AllowedUse == detectlicense.Forbidden {
				return fmt.Errorf("license '%s' is forbidden", license.Name)
			}

			if license.AllowedUse < allowedLicenses {
				return fmt.Errorf("license '%s' should not be used since it doesn't meet the useage requirements", license.Name)
			}
		}
	}
	return nil
}
