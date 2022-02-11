######################################################################
# builder for Js scanning
######################################################################
ARG NODE_VERSION="10"
FROM node:${NODE_VERSION}-alpine as npm_dependency_scanner

RUN apk --no-cache add \
    bash \
    gawk \
    jq

RUN npm install -g license-checker@25.0.1

WORKDIR /scripts
COPY js-mkopensource *.sh customLicenseFormat.json ./
RUN chmod +x *.sh js-mkopensource

WORKDIR /app
COPY npm_dependencies.tar ./
RUN tar xf npm_dependencies.tar && rm -f npm_dependencies.tar

