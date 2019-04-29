FROM busybox:1.30

LABEL maintainer="Steven Xie <dev@stevenxie.me>"
LABEL org.label-schema.schema-version="1.0"
LABEL org.label-schema.name="stevenxie/dockerfiles/frps"
LABEL org.label-schema.vcs-url="https://github.com/stevenxie/dockerfiles"

ARG VERSION=0.26.0

# Copy frp binaries to /usr/bin/.
COPY frps frpc /usr/bin/
COPY frpc.ini frps.ini /etc/frp/

# Expose default FRPS bind port (7000) and FRPS dashboard port (7500).
EXPOSE 7000
EXPOSE 7500

WORKDIR /etc/frp
ENTRYPOINT ["frps"]
CMD ["-c", "frps.ini"]
