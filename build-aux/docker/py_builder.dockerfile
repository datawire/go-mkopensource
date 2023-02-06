########################################
# Python dependency scanner
########################################
ARG PYTHON_IMAGE="need-a-base-image"
FROM golang:1.19-alpine3.15 as builder

WORKDIR /src
COPY . ./

ARG SCRIPTS_HOME
WORKDIR /src/${SCRIPTS_HOME}/cmd/py-mkopensource

RUN mkdir /out
RUN GOOS=linux GARCH=amd64 CGO_ENABLED=0 go build -o /out/ .
WORKDIR /src/${SCRIPTS_HOME}/build-aux/docker/
RUN cp scan-py.sh imports.sh python_dependencies.tar /out/

FROM ${PYTHON_IMAGE} as python_dependency_scanner

ARG APPLICATION_TYPE
ENV APPLICATION_TYPE="${APPLICATION_TYPE}"

RUN apk --no-cache add \
    bash \
    gawk \
    jq

WORKDIR /scripts
COPY --from=builder /out/* ./
RUN chmod +x *.sh py-mkopensource

WORKDIR /app
RUN tar xf /scripts/python_dependencies.tar
