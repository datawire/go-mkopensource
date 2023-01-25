## Command line tool Docs
* [go-mkopensource](/cmd/go-mkopensource/README.md)
* [js-mkopensource](/cmd/js-mkopensource/README.md)

## Building

You may use `go get github.com/datawire/go-mkopensource`, clone the
repo and run `go build .`, or any of the other usual ways of building
a Go program; there is nothing special about `go-mkopensource`.

## Using as a library

The [`github.com/datawire/go-mkopensource/pkg/detectlicense`][detectlicense]
package is good at detecting the licenses in a file

[detectlicense]: https://pkg.go.dev/github.com/datawire/go-mkopensource/pkg/detectlicense

## Design

There are many existing packages to do license detection, such as
[go-license-detector][] or GitHub's [licensee][].  The reason these
are not used is that they are meant to be _informative_, they provide
"best effort" identification of the license.

`go-mkopensource` isn't meant to just be _informative_, it is meant to
be used for _compliance_, if it has any reason at all to be even a
little skeptical of a result, rather than returning its best guess, it
blows up in your face, asking a human to verify the result.

[go-license-detector]: https://github.com/go-enry/go-license-detector
[licensee]: https://github.com/licensee/licensee

## License scanning scripts

Folder `/build-aux` contains scripts to scan licenses for Go, Python 
and Node.Js. Script will generate both `DEPENDENCY_LICENSES.md` and 
`DEPENDENCIES.md`

The following environment variables are used to configure the 
application behaviour.

* `APPLICATION`: Required. Name of the application being scanned.
  It's used in the header of the license files.

* `APPLICATION_TYPE`: Required. Where will the application being 
  scanned run.    
  `internal` is used for anything running on Ambassador Labs servers, 
  and `external` for anything that's deployed to customer machines. 

* `BUILD_HOME` Required. Location of the root folder of the repo to 
  scan.

* `BUILD_TMP`: Required. Folder to use for storing temporary files.

* `GIT_TOKEN` Required. Git token with permissions to pull 
  repositories

* `GO_IMAGE` Required. Image to use for generating Go
  dependencies.

* `PYTHON_PACKAGES`: Optional. List of requirement.txt files to scan.
  Paths should be relative to `BUILD_HOME`.      
  Example:

  `export PYTHON_PACKAGES="./python/requirements.txt ./builder/requirements.txt"`

* `PYTHON_IMAGE`: Required. Image to use for generating Python 
  dependencies.

* `NPM_PACKAGES`: Optional. List of package.json and package-lock.json 
  files to scan. Paths should be relative to `BUILD_HOME`.  
  Example:

  `export NPM_PACKAGES="./tools/sandbox/grpc_web/package.json ./tools/sandbox/grpc_web/package-lock.json"`

* `EXCLUDED_PKG`: Optional. Semicolon separated list of npm packages names that we want to exclude for the validation.
  *Important*: it will restrict the output to the packages (package@version) from being reported in DEPENDENCIES.md and DEPENDENCY_LICENSES.md,
  before to use it, confirm if it is absolutely necessary.
  
  Example:

  `export EXCLUDED_PKG="intro.js@5.0.0;internal-2"`

* `NODE_IMAGE`: Required when `NPM_PACKAGES` is defined. Version 
  of Node.JS to use when running npm dependency scan. Only valid
  version numbers (X.Y.Z) are allowed.  
  Example:

  `NODE_IMAGE=node:14.13.1-alpine`

* `SCRIPTS_HOME`: Required. Location where `go-mkopensource` repo is 
  checked out, relative to  `BUILD_HOME`

To update license information files, set the environment variables 
described above and run `build-aux/generate.sh`

## When scanning fails

The scanner will sometimes fail to detect what kind of licenses a package is using because there's no real standard
stipulating how such files must be organized, and what they must contain. When this happens, and when the failing
package cannot be modified, the scanner must be updated to accommodate the failure. In essence, code must be added
to make the scanner succeed.

Fixing the scanner might take some time, because it's not always self-evident how to do that, or even if it should be
done (the culprit package's owner might be persuaded to fix the problem on their side). Regardless of how the problem
is fixed, it risks blocking progress in projects that use the scanner to produce valid `DEPENDENCIES.md` file.

Blocking progress is never good, and for situations like this, the scanner offers an escape hatch. The `--unparsable-packages`
option. The flag argument is a file name that maps package names to valid SPDX licenses.

Example:
The scanner complains about some package that cannot be parsed:

```
fatal: 2 license-detection errors:
 1. package "sigs.k8s.io/json": could not identify license in file "sigs.k8s.io/json/LICENSE"
 2. package "sigs.k8s.io/json/internal/golang/encoding/json": could not identify license in file "sigs.k8s.io/json/LICENSE"
```

A quick look at the package reveals that it uses an Apache License, but adds extra text at the top of the actual LICENSE
file indicating that it also uses files from golang/encoding/json. We know that golang uses a 3-clause BSD license.  So we consult the [SPDX License List](https://spdx.org/licenses/) to get the canonical
identifiers for the licenses, and add them to an `unparsable-packages.yaml` file to our build system
with the following contents:

```
sigs.k8s.io/json:
  - Apache-2.0
sigs.k8s.io/json/internal/golang/encoding/json:
  - BSD-3-Clause
```

We then use the flag `--unparsable-packages unparsable-packages.yaml` when running `go-mkopensource`.

Example:
In previous versions of this scanner, sometimes the scanner complains about missing dependencies when the scanner gets
the list of all the packages in the file "vendor/modules.txt" using the command "go mod vendor" You can see that in the following output.
```bash
#26 18.21 go: downloading github.com/docker/go-metrics v0.0.0-20180209012529-399ea8c73916
#26 18.24 go: downloading github.com/containerd/cgroups v0.0.0-20200531161412-0dbf7f05ba59
#26 18.28 go: downloading github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da
#26 36.85 github.com/datawire/saas_app/internal/pkg/kubernetes imports
#26 36.85 	k8s.io/client-go/rest imports
#26 36.85 	k8s.io/apimachinery/pkg/util/clock: no required module provides package k8s.io/apimachinery/pkg/util/clock; to add it:
#26 36.85 	go get k8s.io/apimachinery/pkg/util/clock
#26 36.85 /scripts/go-mkopensource: fatal: ["go" "mod" "vendor"]: exit status 1
#26 ERROR: executor failed running [/bin/sh -c /scripts/scan-go.sh]: exit code: 1
```
Now the scanner is smart enough to follow the indications of the "go mod vendor"  install the dependencies, and then 
get the list of packages from the file ''vendor/modules.txt"

Sometimes it isn't possible to install the dependecies sugested by the "go mod vendor" command. 
The scanner will complain with the message "Error installing dependency". In this case the project will require human intervention to solve the problem.

### Remember to always create a ticket!
When a problem arise, remember to always create a ticket so that the problem can be fixed. This will help all users
of the `go-mkopensource` tool and in many cases also make the owner of the failing component aware of the problem.
