ARG VERSION=latest
FROM golang:${VERSION}

# Add labels.
LABEL maintainer="Steven Xie <dev@stevenxie.me>"
LABEL org.label-schema.schema-version="1.0"
LABEL org.label-schema.name=stevenxie/dockerfiles/revive
LABEL org.label-schema.vcs-url="https://github.com/stevenxie/dockerfiles"

# Copy files.
COPY ./revive ./goimports /bin/
COPY ./lint.sh /scripts/

ENTRYPOINT ["sh"]
CMD ["/scripts/lint.sh"]
