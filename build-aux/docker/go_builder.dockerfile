########################################
# Go dependency scanner
########################################
ARG GO_IMAGE="base-image-unknown"
FROM golang:1.19-alpine3.15 as builder

WORKDIR /src
COPY . ./

ARG SCRIPTS_HOME
WORKDIR /src/${SCRIPTS_HOME}/cmd/go-mkopensource

RUN mkdir /out
RUN GOOS=linux GARCH=amd64 CGO_ENABLED=0 go build -o /out/ .
WORKDIR /src/${SCRIPTS_HOME}/build-aux/docker/
RUN cp scan-go.sh imports.sh /out/

FROM ${GO_IMAGE} as go_dependency_scanner

ARG UNPARSABLE_PACKAGE
ARG APPLICATION_TYPE
ENV APPLICATION_TYPE="${APPLICATION_TYPE}"

RUN apk --no-cache add \
    bash \
    curl \
    gawk \
    git \
    jq

WORKDIR /data
RUN set -ex; GO_VERSION=$(go version | sed -E 's/.*go([1-9\.]*).*/\1/') && \
    GO_TAR="go${GO_VERSION}.src.tar.gz" && \
    curl -o "${GO_TAR}" --fail -L "https://dl.google.com/go/go${GO_VERSION}.src.tar.gz"

WORKDIR /app
COPY . ./

WORKDIR /scripts
COPY --from=builder /out/* ./
RUN chmod +x *.sh go-mkopensource

ARG GIT_TOKEN
RUN git config --global url."https://$GIT_TOKEN:@github.com/".insteadOf "https://github.com/"

WORKDIR /app
RUN if [[ -z "$UNPARSABLE_PACKAGE" ]] ; then /scripts/scan-go.sh; else /scripts/scan-go.sh --unparsable-packages $UNPARSABLE_PACKAGE ; fi

FROM scratch as license_output
COPY --from=go_dependency_scanner /temp/* /
