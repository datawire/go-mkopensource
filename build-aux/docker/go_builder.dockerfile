########################################
# Go dependency scanner
########################################
ARG GO_BUILDER="base-image-unknown"
FROM ${GO_BUILDER} as go_dependency_scanner

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

ARG SCRIPTS_HOME
RUN ln -s $(realpath "${SCRIPTS_HOME}/build-aux/docker/") /scripts
RUN chmod +x /scripts/*.sh /scripts/go-mkopensource

ARG GIT_TOKEN
RUN git config --global url."https://$GIT_TOKEN:@github.com/".insteadOf "https://github.com/"

WORKDIR /app
RUN /scripts/scan-go.sh

FROM scratch as license_output
COPY --from=go_dependency_scanner /temp/* /
