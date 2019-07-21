ARG VERSION=v0.10.0
FROM gcr.io/kaniko-project/executor:debug-${VERSION}

# Add labels.
LABEL maintainer="Steven Xie <dev@stevenxie.me>"
LABEL org.label-schema.schema-version="1.0"
LABEL org.label-schema.name=stevenxie/dockerfiles/kaniko-drone
LABEL org.label-schema.vcs-url="https://github.com/stevenxie/dockerfiles"

# Copy into image as /kaniko/runner.
COPY kaniko-drone /kaniko/runner
ENTRYPOINT ["/kaniko/runner"]
