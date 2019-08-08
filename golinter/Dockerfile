FROM golang:1-alpine3.10

# Add labels.
LABEL maintainer="Steven Xie <dev@stevenxie.me>"
LABEL org.label-schema.schema-version="1.0"
LABEL org.label-schema.name=stevenxie/dockerfiles/revive
LABEL org.label-schema.vcs-url="https://github.com/stevenxie/dockerfiles"

COPY ./revive ./goimports ./lint /bin/
ENTRYPOINT ["lint"]
