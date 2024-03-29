build:  go-mkopensource py-mkopensource js-mkopensource
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

lint: tools/bin/golangci-lint
	tools/bin/golangci-lint run ./...
.PHONY: lint

tools/bin/%: tools/src/%/pin.go tools/src/%/go.mod
	cd $(<D) && GOOS= GOARCH= go build -o $(abspath $@) $$(sed -En 's,^import "(.*)".*,\1,p' pin.go)

.DELETE_ON_ERROR:
.PHONY: FORCE
FORCE:
