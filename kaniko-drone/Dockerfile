ARG VERSION=v0.9.0
FROM gcr.io/kaniko-project/executor:debug-${VERSION}

# Copy into image as /kaniko/runner.
COPY kaniko-drone /kaniko/runner
ENTRYPOINT ["/kaniko/runner"]
