FROM alpine:3.9

LABEL maintainer="Steven Xie <dev@stevenxie.me>"
LABEL org.label-schema.schema-version="1.0"
LABEL org.label-schema.name="stevenxie/dockerfiles/frps"
LABEL org.label-schema.vcs-url="https://github.com/stevenxie/dockerfiles"

ARG VERSION=0.26.0

# Download frp binaries to /usr/bin/.
WORKDIR /tmp
RUN wget -O frp.tar.gz \
      "https://github.com/fatedier/frp/releases/download/v${VERSION}/frp_${VERSION}_linux_amd64.tar.gz" && \
    tar -zxvf frp.tar.gz && \
    cd "frp_${VERSION}_linux_amd64" && \
    mv frps frpc /usr/bin/ && \
    mkdir -p /etc/frp && \
    mv frpc.ini frps.ini /etc/frp/ && \
    cd .. && \
    rm -rf *

# Expose default FRPS bind port (7000) and FRPS dashboard port (7500).
EXPOSE 7000
EXPOSE 7500

WORKDIR /etc/frp
ENTRYPOINT ["frps"]
CMD ["-c", "frps.ini"]
