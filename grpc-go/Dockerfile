ARG GOVERSION=latest

FROM golang:${GOVERSION}

# Add labels.
LABEL maintainer="Steven Xie <dev@stevenxie.me>"
LABEL org.label-schema.schema-version="1.0"
LABEL org.label-schema.name=stevenxie/dockerfiles/grpc-go
LABEL org.label-schema.vcs-url="https://github.com/stevenxie/dockerfiles"

COPY protoc protoc-gen-go /usr/local/bin/
