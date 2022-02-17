package dependency

import . "github.com/datawire/go-mkopensource/pkg/detectlicense"

//nolint:gochecknoglobals // Would be 'const'.
var hardcodedJsDependencies = map[string][]string{
	"cyclist@0.2.2":                {MIT.Name},
	"doctrine@1.5.0":               {BSD2.Name, Apache2.Name},
	"emitter-component@1.1.1":      {MIT.Name},
	"flexboxgrid@6.3.1":            {Apache2.Name},
	"indexof@0.0.1":                {MIT.Name},
	"intro.js@4.1.0":               {AGPL3OrLater.Name},
	"node-forge@0.10.0":            {BSD3.Name},
	"pako@1.0.10":                  {MIT.Name},
	"regenerator-transform@0.10.1": {BSD2.Name},
	"regjsparser@0.1.5":            {BSD2.Name},
}
