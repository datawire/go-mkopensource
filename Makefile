build: go-mkopensource
.PHONY: build

go-mkopensource: FORCE
	cd cmd/go-mkopensource; \
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


#############################################################
## Generate license information
#############################################################

TOKEN=$(shell if [ -n "$${GIT_TOKEN}" ]; then echo "$${GIT_TOKEN}"; else grep oauth_token ~/.config/hub | awk '{print $$2}'; fi)

LICENSE_TMP_DIR=./licenses.tmp

clean-dependency-info:
	rm -fR "$(CURDIR)/$(LICENSE_TMP_DIR)"
.PHONY: clean-dependency-info

generate-dependency-info:
	set -e; { \
		if [ ! "$(TOKEN)" ]; then \
			printf '>>> $(bold)$(ccgreen)No GIT_TOKEN provided. Make sure you have hub installed and configured. See: https://github.com/github/hub$(sgr0)\n'; \
			exit 1; \
		fi; \
		\
		export APPLICATION="Ambassador Cloud"; \
		export APPLICATION_TYPE="internal"; \
		export BUILD_HOME='.'; \
		export BUILD_TMP="$(LICENSE_TMP_DIR)/output"; \
		export SCRIPTS_HOME="."; \
		export GO_IMAGE='golang:1.17-alpine3.15' \
		export GIT_TOKEN="$(TOKEN)"; \
		"./build-aux/generate.sh"; \
	}
.PHONY: generate-dependency-info
