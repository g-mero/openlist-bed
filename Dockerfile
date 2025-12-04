FROM golang:1.25.3-bookworm AS builder

# RUN sed -i 's/deb.debian.org/mirrors.ustc.edu.cn/g' /etc/apt/sources.list.d/debian.sources
# RUN sed -i 's/security.debian.org/mirrors.ustc.edu.cn/g' /etc/apt/sources.list.d/debian.sources

# Add Debian sid repository for latest libvips
RUN echo "deb http://ftp.hk.debian.org/debian sid main" > /etc/apt/sources.list.d/sid.list \
    && apt-get update \
    && apt-get install -y --no-install-recommends -t sid libvips-dev \
    && rm /etc/apt/sources.list.d/sid.list \
    && apt-get update \
    && mkdir /build

COPY go.mod /build
RUN cd /build && go mod tidy

COPY . /build
RUN cd /build \
    && CGO_ENABLED=1 go build -ldflags="-s -w" -o main cmd/server/main.go

FROM debian:bookworm-slim

# RUN sed -i 's/deb.debian.org/mirrors.ustc.edu.cn/g' /etc/apt/sources.list.d/debian.sources
# RUN sed -i 's/security.debian.org/mirrors.ustc.edu.cn/g' /etc/apt/sources.list.d/debian.sources

# Add Debian sid repository for latest libvips runtime
RUN echo "deb http://ftp.hk.debian.org/debian sid main" > /etc/apt/sources.list.d/sid.list \
    && apt-get update \
    && apt-get install -y --no-install-recommends -t sid libvips42t64 \
    && apt-get install -y --no-install-recommends ca-certificates libjemalloc2 libtcmalloc-minimal4 \
    && rm /etc/apt/sources.list.d/sid.list \
    && rm -rf /var/lib/apt/lists/*

COPY --from=builder /build/main /opt/main

WORKDIR /opt
VOLUME /opt/data
EXPOSE 6001
CMD ["/opt/main"]