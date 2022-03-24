######################################################################
# builder for Js scanning
######################################################################
ARG NODE_IMAGE="need-a-base-image"
FROM golang:1.17-alpine3.15 as builder

WORKDIR /src
COPY . ./

ARG SCRIPTS_HOME
WORKDIR /src/${SCRIPTS_HOME}/cmd/js-mkopensource

RUN mkdir /out
RUN GOOS=linux GARCH=amd64 CGO_ENABLED=0 go build -o /out/ .
WORKDIR /src/${SCRIPTS_HOME}/build-aux/docker/
RUN cp scan-js.sh imports.sh customLicenseFormat.json npm_dependencies.tar /out/

FROM ${NODE_IMAGE} as npm_dependency_scanner

ARG APPLICATION_TYPE
ENV APPLICATION_TYPE="${APPLICATION_TYPE}"
ENV EXCLUDED_PKG="${EXCLUDED_PKG}"

RUN apk --no-cache add \
    bash \
    gawk \
    jq

RUN npm install -g license-checker@25.0.1

WORKDIR /scripts
COPY --from=builder /out/* ./
RUN chmod +x *.sh js-mkopensource

WORKDIR /app
RUN tar xf /scripts/npm_dependencies.tar
