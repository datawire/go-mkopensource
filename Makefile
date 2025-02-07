build:  go-mkopensource js-mkopensource
.PHONY: build

%-mkopensource: FORCE cmd/%-mkopensource
	cd cmd/$*-mkopensource; \
	go build .

check:
	go test -race ./...
.PHONY: check

generate:
	go generate ./...
.PHONY: generate

.DELETE_ON_ERROR:
.PHONY: FORCE
FORCE:
