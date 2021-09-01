package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/datawire/go-mkopensource/pkg/mkopensource/golist"
)

func VendorList() ([]golist.Package, error) {
	cmd := exec.Command("go", "mod", "vendor")
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("%q: %w", []string{"go", "mod", "vendor"}, err)
	}

	file, err := os.Open("vendor/modules.txt")
	if err != err {
		return nil, err
	}
	defer file.Close()

	var pkgs []golist.Package
	var curModule *golist.Module
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") {
			parts := strings.Split(line, " ")
			switch len(parts) {
			case 3:
				// ok
			case 5, 6:
				if parts[3] != "=>" {
					return nil, fmt.Errorf("malformed line in vendor/modules.txt: %q", line)
				}
			default:
				return nil, fmt.Errorf("malformed line in vendor/modules.txt: %q", line)
			}
			modname := parts[1]
			modules, err := golist.GoListModules([]string{"-mod=vendor"}, []string{modname})
			if err != nil {
				return nil, err
			}
			if len(modules) != 1 {
				return nil, errors.New("unexpected output from go list")
			}
			curModule = &modules[0]
		} else {
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
