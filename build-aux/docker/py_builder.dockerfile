########################################
# Python dependency scanner
########################################
ARG PYTHON_BUILDER="need-a-base-image"
FROM ${PYTHON_BUILDER} as python_dependency_scanner

RUN apk --no-cache add \
    bash \
    gawk \
    jq

WORKDIR /scripts
COPY py-mkopensource *.sh ./
RUN chmod +x *.sh py-mkopensource

WORKDIR /app
COPY python_dependencies.tar ./
RUN tar xf python_dependencies.tar && rm -f python_dependencies.tar
