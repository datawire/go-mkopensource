package main

type dependencyInfo struct {
	Dependencies []dependency `json:"dependencies"`
}

type dependency struct {
	Name     string   `json:"name"`
	Version  string   `json:"version"`
	Licenses []string `json:"licenses"`
}
