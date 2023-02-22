package dependency

import . "github.com/datawire/go-mkopensource/pkg/detectlicense"

//nolint:gochecknoglobals // Would be 'const'.
var hardcodedJsDependencies = map[string][]License{
	"cyclist@0.2.2":                {MIT},
	"doctrine@1.5.0":               {BSD2, Apache2},
	"emitter-component@1.1.1":      {MIT},
	"flexboxgrid@6.3.1":            {Apache2},
	"indexof@0.0.1":                {MIT},
	"intro.js@4.1.0":               {AGPL3OrLater},
	"json-schema@0.2.3":            {AFL21},
	"node-forge@0.10.0":            {BSD3},
	"pako@1.0.10":                  {MIT},
	"regenerator-transform@0.10.1": {BSD2},
	"regjsparser@0.1.5":            {BSD2},
	"node-forge@1.3.1":             {BSD3},
}
