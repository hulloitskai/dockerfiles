ARG VERSION=latest
FROM alpine/git:${VERSION}

# Add labels.
LABEL maintainer="Steven Xie <dev@stevenxie.me>"
LABEL org.label-schema.schema-version="1.0"
LABEL org.label-schema.name=stevenxie/dockerfiles/precond-drone
LABEL org.label-schema.vcs-url="https://github.com/stevenxie/dockerfiles"

# Copy scripts.
COPY ./precond.sh /scripts/

ENTRYPOINT ["sh"]
CMD ["/scripts/precond.sh"]
