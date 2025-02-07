package detectlicense

import (
	"fmt"
)

// knownDependencies will return a list of licenses for any dependency that has been
// hardcoded due to the difficulty to parse the license file(s).
func knownDependencies(dependencyName string, dependencyVersion string) (licenses []License, ok bool) {
	hardcodedGoDependencies := map[string][]License{
		"github.com/josharian/intern@v1.0.1-0.20211109044230-42b52b674af5":                  {MIT}, // License had a funny filename, fixed in https://github.com/josharian/intern/pull/2
		"github.com/garyburd/redigo/internal@v0.0.0-20150301180006-535138d7bcd7":            {Apache2}, // Just had a note in the README, a LICENSE file wasn't added until 1.0.0
		"github.com/garyburd/redigo/redis@v0.0.0-20150301180006-535138d7bcd7":               {Apache2}, // Just had a note in the README, a LICENSE file wasn't added until 1.0.0
	}

	licenses, ok = hardcodedGoDependencies[fmt.Sprintf("%s@%s", dependencyName, dependencyVersion)]
	return
}
