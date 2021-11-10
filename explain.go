package main

func ExplainError(err error) string {
	switch {

	case err.Error() == `package "github.com/josharian/intern": could not identify a license for all sources (had no global LICENSE file)`:

		return wordwrap(4, 72, `For github.com/josharian/intern in particular, this probably
			means that you are depending on an old version; upgrading to intern
			v1.0.1-0.20211109044230-42b52b674af5 or later should resolve this.`)

	default:

		return wordwrap(4, 72, `This probably means that you added or upgraded a dependency,
			and the automated opensource-license-checker can't confidently detect what
			the license is.  (This is a good thing, because it is reminding you to check
			the license of libraries before using them.)

			You need to update the
			"github.com/datawire/go-mkopensource/pkg/detectlicense/licenses.go" file to
			correctly detect the license.`)

	}
}
