package main

type DependencyInfo struct {
	Dependencies []dependency      `json:"dependencies"`
	Licenses     map[string]string `json:"licenseInfo"`
}

type dependency struct {
	Name     string   `json:"name"`
	Version  string   `json:"version"`
	Licenses []string `json:"licenses"`
}

func NewDependencyInfo() DependencyInfo {
	return DependencyInfo{
		Dependencies: []dependency{},
		Licenses:     map[string]string{},
	}
}
