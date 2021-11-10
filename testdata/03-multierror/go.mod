module testmod

go 1.17

require (
	example.com/apache-patent v0.0.0-00010101000000-000000000000
	example.com/cc-sa v0.0.0-00010101000000-000000000000
	example.com/gpl v0.0.0-00010101000000-000000000000
	github.com/josharian/intern v1.0.0
)

replace (
	example.com/apache-patent => ./apache-patent
	example.com/cc-sa => ./cc-sa
	example.com/gpl => ./gpl
)
