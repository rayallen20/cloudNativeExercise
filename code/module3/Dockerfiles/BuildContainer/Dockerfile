FROM ubuntu

# so apt-get doesn't complain
ENV DEBIAN_FRONTEND=noninteractive
RUN sed -i 's/^exit 101/exit 0' /usr/sbin/policy-rc.d

RUN \
    apt-get update && \
    apt-get install -y ca-certificates && \
    apt-get install -y curl && \
    rm -rf /var/lib/apt/lists/*

ADD ./bin/myapp myapp
ENTRYPOINT ["/myapp"]