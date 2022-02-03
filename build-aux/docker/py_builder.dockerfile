########################################
# Python dependency scanner
########################################

FROM docker.io/frolvlad/alpine-glibc:alpine-3.15 as python_dependency_scanner

WORKDIR /buildroot

ENV PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/usr/local/go/bin:/buildroot/bin

ARG PYTHON_VERSION="~3.9.7"
RUN apk --no-cache add \
    bash \
    gcc \
    make \
    musl-dev \
    curl \
    cython \
    docker-cli \
    git \
    iptables \
    jq \
    libcap \
    libcap-dev \
    libffi-dev \
    ncurses \
    openssl-dev \
    py3-pip=~20.3.4 \
    python3=$PYTHON_VERSION \
    python3-dev \
    rust \
    cargo \
    patchelf \
    rsync \
    sudo \
    yaml-dev \
    && ln -s /usr/bin/python3 /usr/bin/python \
    && chmod u+s $(which docker)

# Consult
# https://github.com/jazzband/pip-tools/#versions-and-compatibility to
# select a pip-tools version that corresponds to the 'py3-pip' and
# 'python3' versions above.
RUN pip3 install pip-tools==6.3.1

# The YAML parser is... special. To get the C version, we need to install Cython and libyaml, then
# build it locally -- just using pip won't work.
#
# Download, build, and install PyYAML.
RUN mkdir /tmp/pyyaml && \
  cd /tmp/pyyaml && \
  curl -o pyyaml-5.4.1.1.tar.gz -L https://github.com/yaml/pyyaml/archive/refs/tags/5.4.1.1.tar.gz && \
  tar xzf pyyaml-5.4.1.1.tar.gz && \
  cd pyyaml-5.4.1.1 && \
  python3 setup.py --with-libyaml install

# orjson is also special.  The wheels on PyPI rely on glibc, so we
# need to use cargo/rustc/patchelf to build a musl-compatible version.
RUN pip3 install orjson==3.6.6

WORKDIR /scripts
COPY py-mkopensource *.sh ./
RUN chmod +x *.sh py-mkopensource

RUN ln -s /buildroot /app
WORKDIR /app
COPY python_dependencies.tar ./
RUN tar xf python_dependencies.tar && rm -f python_dependencies.tar
