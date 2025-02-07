########################################
# Go dependency scanner
########################################
ARG GO_IMAGE="base-image-unknown"
FROM golang:1.23.6-alpine3.21 AS builder

ENV GOCACHE=/root/.cache/go-build
RUN mkdir -p "${GOCACHE}"

ENV GOMODCACHE=/root/go/pkg/mod
RUN mkdir -p "${GOMODCACHE}"

WORKDIR /src
COPY . ./
ARG SCRIPTS_HOME
WORKDIR /src/${SCRIPTS_HOME}/cmd/go-mkopensource

RUN mkdir /out

RUN --mount=type=cache,target=/root/.cache/go-build --mount=type=cache,target=/root/go/pkg/mod \
    GOOS=linux GARCH=amd64 CGO_ENABLED=0 go build -o /out/ .

WORKDIR /src/${SCRIPTS_HOME}/build-aux/docker/
RUN cp scan-go.sh imports.sh /out/

FROM ${GO_IMAGE} AS go_dependency_scanner

ENV GOCACHE=/root/.cache/go-build
RUN mkdir -p "${GOCACHE}"

ENV GOMODCACHE=/root/go/pkg/mod
RUN mkdir -p "${GOMODCACHE}"

ARG UNPARSABLE_PACKAGE
ENV UNPARSABLE_PACKAGE="${UNPARSABLE_PACKAGE}"
ARG PROPRIETARY_PACKAGES
ENV PROPRIETARY_PACKAGES="${PROPRIETARY_PACKAGES}"
ARG APPLICATION_TYPE
ENV APPLICATION_TYPE="${APPLICATION_TYPE}"

RUN --mount=type=cache,target=/var/cache/apk,sharing=locked \
    apk add \
    bash \
    curl \
    gawk \
    git \
    jq

WORKDIR /data
RUN set -ex; GO_VERSION=$(go version | sed -E 's/.*go([0-9\.]*).*/\1/') && \
    GO_TAR="go${GO_VERSION}.src.tar.gz" && \
    curl -o "${GO_TAR}" --fail -L "https://dl.google.com/go/go${GO_VERSION}.src.tar.gz"

ARG GIT_TOKEN
RUN git config --global url."https://$GIT_TOKEN:@github.com/".insteadOf "https://github.com/"

WORKDIR /scripts
COPY --from=builder /out/* ./
RUN chmod +x *.sh go-mkopensource

WORKDIR /app
COPY . ./

RUN test -z "${UNPARSABLE_PACKAGE}" || stat "${UNPARSABLE_PACKAGE}" > /dev/null
RUN test -z "${PROPRIETARY_PACKAGES}" || stat "${PROPRIETARY_PACKAGES}" > /dev/null

RUN --mount=type=cache,target=/root/.cache/go-build --mount=type=cache,target=/root/go/pkg/mod \
    go mod download

RUN --mount=type=cache,target=/root/.cache/go-build --mount=type=cache,target=/root/go/pkg/mod \
    /scripts/scan-go.sh

RUN cp go.mod go.sum /temp/

FROM scratch AS license_output
COPY --from=go_dependency_scanner /temp/* /
COPY --from=go_dependency_scanner /app/go.mod /app/go.sum /
