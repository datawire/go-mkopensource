###
# This dockerfile builds the base image for the builder container. See
# the main Dockerfile for more information about what the builder
# container is and how code in this repo is built.
#
# Originally this base was built as part of the builder container's
# bootstrap process. We discovered that minor network interruptions
# would break these steps, and such interruptions were common on our
# cloud CI system. We decided to separate out these steps so that any
# one of them is much less likely to be the cause of a network-related
# failure, i.e. a flake.
#
# See the comment before the build_builder_base() function in builder.sh
# to see when and how often this base image is built and pushed.
##

# This argument controls the base image that is used for our build
# container.
########################################
# Python dependency scanner
########################################
FROM docker.io/frolvlad/alpine-glibc:alpine-3.12_glibc-2.32 as python_dependency_scanner

WORKDIR /buildroot

ENV PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/usr/local/go/bin:/buildroot/bin

ARG PYTHON_VERSION="~3.8.10"
RUN apk --no-cache add \
    bash \
    gcc \
    make \
    musl-dev \
    curl \
    cython \
    gawk \
    jq \
    libcap \
    libcap-dev \
    libffi-dev \
    ncurses \
    openssh-client \
    openssl-dev \
    py3-pip \
    python3=$PYTHON_VERSION \
    python3-dev \
    rsync \
    sudo \
    && ln -s /usr/bin/python3 /usr/bin/python

# We _must_ pin pip to a version before 20.3 because orjson appears to only have
# PEP513 compatible wheels, which are supported before 20.3 but (apparently)
# not in 20.3. We can only upgrade pip to 20.3 after we verify that orjson has
# PEP600 compatible wheels for our linux platform, or we start building orjson
# from source using a rust toolchain.
RUN pip3 install -U pip==20.2.4 pip-tools==5.3.1

WORKDIR /scripts
COPY py-mkopensource *.sh ./
RUN chmod +x *.sh py-mkopensource

RUN ln -s /buildroot /app

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

