## Command line tool Docs
* [go-mkopensource](/cmd/go-mkopensource/README.md)

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
little skeptical of a result, rather than returnit its best guess, it
blows up in your face, asking a human to verify the result.

[go-license-detector]: https://github.com/go-enry/go-license-detector
[licensee]: https://github.com/licensee/licensee
