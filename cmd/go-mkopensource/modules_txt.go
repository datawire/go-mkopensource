package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"regexp"
	"runtime"
	"strings"

	"github.com/datawire/go-mkopensource/pkg/golist"
)

const DEFAULT_DEPENDENCY_REGEX = "default"

type Modules struct {
	removedDependenciesRegex map[string]*regexp.Regexp
}

func NewModules() (m *Modules) {
	m = &Modules{
		removedDependenciesRegex: map[string]*regexp.Regexp{},
	}

	m.removedDependenciesRegex["go1.16"] = regexp.MustCompile(`\t+([^\s]+): module .* found .* but does not contain package`)
	m.removedDependenciesRegex["go1.17"] = m.removedDependenciesRegex["go1.16"]
	m.removedDependenciesRegex["go1.18"] = regexp.MustCompile(`\t+([^\s]+): no required module provides package`)
	m.removedDependenciesRegex["go1.19"] = m.removedDependenciesRegex["go1.18"]
	m.removedDependenciesRegex[DEFAULT_DEPENDENCY_REGEX] = m.removedDependenciesRegex["go1.18"]

	return m
}

// VendorList returns a listing of all packages in
// `vendor/modules.txt`, which is superior to `go list -deps` in that
// it includes dependencies for all platforms and build
// configurations, but inferior in that it cannot be asked to only
// consider dependencies of a specific package rather than the whole
// module.
func (m *Modules) VendorList() ([]golist.Package, error) {
	// References: In the Go stdlib source code, see
	// - `cmd/go/internal/modcmd/vendor.go` for the code that writes modules.txt, and
	// - `cmd/go/internal/modload/vendor.go` for the code that parses it.
	cmd := exec.Command("go", "mod", "vendor")
	if out, err := cmd.CombinedOutput(); err != nil {
		if errInstall := m.findAndGetDependencies(string(out)); errInstall != nil {
			if err := m.tryRemoveUnavailableDependencies(); err != nil {
				return nil, err
			}
		}

		// Run go mod vendor again to update vendored dependencies
		cmd = exec.Command("go", "mod", "vendor")
		if err := cmd.Run(); err != nil {
			return nil, err
		}
	}

	file, err := os.Open("vendor/modules.txt")
	if err != nil {
		if os.IsNotExist(err) {
			// If there are no dependencies outside of stdlib.
			return nil, nil
		}
		return nil, err
	}
	defer file.Close()

	var pkgs []golist.Package
	var curModuleName string
	var curModule *golist.Module // lazily populated from curModuleName
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "## ") {
			// These lines are introduced in Go 1.17 and indicate (1) the Go version in
			// go.mod, and (2) whether we implicitly or explicitly depend on it; neither
			// of which are things we care about.
		} else if strings.HasPrefix(line, "# ") {
			parts := strings.Split(line, " ")
			// Just do some quick validation of the line format.  We're not tring to be
			// super strict with the validation, just a quick check that we're not
			// looking at something totally insane.
			switch len(parts) {
			case 3:
				// 0 1      2
				// # module version
			case 4, 5, 6:
				// 0 1      2       3      4       5
				// # module version =>     module version
				// # module =>      module version
				// # module version =>     path
				// # module =>      path
				if parts[2] != "=>" && parts[3] != "=>" {
					return nil, fmt.Errorf("malformed line in vendor/modules.txt: %q", line)
				}
			default:
				return nil, fmt.Errorf("malformed line in vendor/modules.txt: %q", line)
			}
			// Defer looking up curModule from curModuleName until we actually need it;
			// a non-used replaced module might not be present in `vendor/`.  We could
			// instead download it by using `-mod=readonly` instead of `-mod=vendor`,
			// but what would the point in that be?
			curModuleName = parts[1]
			curModule = nil
		} else {
			if curModule == nil && curModuleName != "" {
				modules, err := golist.GoListModules([]string{"-mod=vendor"}, []string{curModuleName})
				if err != nil {
					return nil, err
				}
				if len(modules) != 1 {
					return nil, errors.New("unexpected output from go list")
				}
				curModule = &modules[0]
			}
			pkgname := line
			pkgs = append(pkgs, golist.Package{
				Dir:        "vendor/" + pkgname,
				ImportPath: pkgname,
				Name:       path.Base(pkgname),
				Module:     curModule,
				DepOnly:    true,
			})
		}
	}

	return pkgs, nil
}

func (m *Modules) tryRemoveUnavailableDependencies() error {
	dependenciesLeftToRemove := 0
	for {
		cmd := exec.Command("go", "mod", "vendor")
		out, err := cmd.CombinedOutput()
		if err == nil {
			break
		}

		lines := strings.Split(string(out), "\n")
		removedDependencies, err := m.getRemovedDependencies(lines)
		if err != nil {
			return err
		}

		if len(removedDependencies) == 0 {
			return fmt.Errorf("%q: %w", []string{"go", "mod", "vendor"}, err)
		}

		if len(removedDependencies) == dependenciesLeftToRemove {
			return fmt.Errorf("number of dependencies to remove didn't change, so removal is not working as expected")
		}

		err = m.updateDependenciesOfRemovedPackage(removedDependencies[0])
		if err != nil {
			return fmt.Errorf("error updating removed dependencies: %w", err)
		}

		dependenciesLeftToRemove = len(removedDependencies)
	}
	return nil
}

func (m *Modules) findAndGetDependencies(outputFromModVendor string) error {
	lines := strings.Split(outputFromModVendor, "\n")
	var dependenciesToInstall []string
	for _, line := range lines {
		if strings.Contains(line, "go get") {
			dependenciesToInstall = append(dependenciesToInstall, line)
		}
	}
	if len(dependenciesToInstall) <= 0 {
		log.Println(outputFromModVendor)
		return fmt.Errorf("none dependency required installation")
	}
	for _, dependency := range dependenciesToInstall {
		command := strings.Split(strings.TrimSpace(dependency), " ")
		cmd := exec.Command(command[0], command[1:]...)
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			log.Printf("Error installing dependency %v", err)
			return fmt.Errorf("%q: %w", []string{"go", "mod", "vendor"}, err)
		}
	}
	return nil
}

func (m *Modules) getRemovedDependencies(lines []string) (removedDependencies []string, err error) {
	re, err := m.getRemovedDependencyRegex()
	if err != nil {
		return nil, err
	}

	for _, line := range lines {
		p := re.FindStringSubmatch(line)
		if len(p) == 2 {
			dependency := p[1]
			log.Printf("found package %s to remove\n", dependency)
			removedDependencies = append(removedDependencies, dependency)
		}
	}
	return removedDependencies, nil
}

func (m *Modules) getRemovedDependencyRegex() (*regexp.Regexp, error) {
	goVersion, err := m.getGoSemVer()
	if err != nil {
		return nil, err
	}

	if re, ok := m.removedDependenciesRegex[goVersion]; ok {
		return re, nil
	}

	return m.removedDependenciesRegex[DEFAULT_DEPENDENCY_REGEX], nil
}

func (m *Modules) getGoSemVer() (version string, err error) {
	goVersionRegex := regexp.MustCompile(`^(go[0-9]+\.[0-9]+)\.[0-9]+$`)
	runtimeVersion := runtime.Version()

	match := goVersionRegex.FindStringSubmatch(runtimeVersion)
	if match == nil {
		return "", fmt.Errorf("could not get go version from %s", runtimeVersion)
	}

	return match[1], nil
}

func (m *Modules) updateDependenciesOfRemovedPackage(removedPackage string) (err error) {
	whyCmd := exec.Command("go", "mod", "why", removedPackage)
	out, err := whyCmd.Output()
	if err != nil {
		log.Printf("'go mod why' failed:\n%s\n", err.(*exec.ExitError).Stderr)
		return fmt.Errorf("'go mod why' failed with error %w", err)
	}

	outputLines := strings.Split(strings.TrimSpace(string(out)), "\n")
	if len(outputLines) <= 3 {
		return fmt.Errorf("Package %s is being imported directly and can't be removed", removedPackage)
	}

	for i := 2; i <= len(outputLines)-2; i++ {
		dependencyToUpdate := outputLines[i]
		updateCmd := exec.Command("go", "get", dependencyToUpdate)
		out, err = updateCmd.CombinedOutput()
		if err != nil {
			log.Printf("'go get %s' failed:\n%s\n", dependencyToUpdate, out)
			return fmt.Errorf("'go get %s' failed with error %w\n", dependencyToUpdate, err)
		}

		tidyCmd := exec.Command("go", "mod", "tidy")
		out, err = tidyCmd.CombinedOutput()
		if err != nil {
			log.Printf("'go mod tidy' failed:\n%s\n", out)
			return fmt.Errorf("'go mod tidy' failed with error %w", err)
		}
	}

	return nil
}
