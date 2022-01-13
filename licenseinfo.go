package main

type dependencyInfo struct {
	Dependencies []dependency      `json:"dependencies"`
	Licenses     map[string]string `json:"licenseInfo"`
}

type dependency struct {
	Name     string   `json:"name"`
	Version  string   `json:"version"`
	Licenses []string `json:"licenses"`
}
