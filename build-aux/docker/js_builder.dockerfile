######################################################################
# builder for Js scanning
######################################################################
ARG NODE_IMAGE="need-a-base-image"
FROM golang:1.19-alpine3.15 as builder

ENV GOCACHE=/root/.cache/go-build
RUN mkdir -p "${GOCACHE}"

ENV GOMODCACHE=/root/go/pkg/mod
RUN mkdir -p "${GOMODCACHE}"

WORKDIR /src
COPY . ./

ARG SCRIPTS_HOME
WORKDIR /src/${SCRIPTS_HOME}/cmd/js-mkopensource

RUN mkdir /out
RUN --mount=type=cache,target=/root/.cache/go-build --mount=type=cache,target=/root/go/pkg/mod \
    GOOS=linux GARCH=amd64 CGO_ENABLED=0 go build -o /out/ .

WORKDIR /src/${SCRIPTS_HOME}/build-aux/docker/
RUN cp scan-js.sh imports.sh customLicenseFormat.json npm_dependencies.tar /out/

FROM ${NODE_IMAGE} as npm_dependency_scanner

ARG APPLICATION
ENV APPLICATION="${APPLICATION}"
ARG APPLICATION_TYPE
ENV APPLICATION_TYPE="${APPLICATION_TYPE}"
ARG EXCLUDED_PKG
ENV EXCLUDED_PKG="${EXCLUDED_PKG}"
ARG USER_ID
ENV USER_ID="${USER_ID}"

RUN --mount=type=cache,target=/var/cache/apk,sharing=locked \
    apk add \
    bash \
    gawk \
    jq

WORKDIR /scripts
COPY --from=builder /out/* ./
RUN chmod +x *.sh js-mkopensource

WORKDIR /app
RUN tar xf /scripts/npm_dependencies.tar
RUN --mount=type=cache,target=/root/.npm,sharing=locked \
    npm set cache /root/.npm && \
    npm install -g license-checker@25.0.1

RUN --mount=type=cache,target=/root/.npm,sharing=locked \
    /scripts/scan-js.sh

FROM scratch as license_output
COPY --from=npm_dependency_scanner /temp/js_dependencies.txt /temp/js_licenses.txt /
