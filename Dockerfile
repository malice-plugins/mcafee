FROM ubuntu:xenial

LABEL maintainer "https://github.com/blacktop"

LABEL malice.plugin.repository = "https://github.com/malice-plugins/mcafee.git"
LABEL malice.plugin.category="av"
LABEL malice.plugin.mime="*"
LABEL malice.plugin.docker.engine="*"

# Install McAfee AV
RUN set -x \
    && apt-get update \
    && apt-get install -yq ca-certificates curl wget unzip --no-install-recommends \
    && echo "===> Install McAfee..." \
    && mkdir -p /usr/local/uvscan \
    && curl http://b2b-download.mcafee.com/products/evaluation/vcl/l64/vscl-l64-604-e.tar.gz \
    | tar -xzf - -C /usr/local/uvscan \
    && echo "===> Clean up unnecessary files..." \
    && apt-get purge -y --auto-remove curl $(apt-mark showauto) \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/* /var/cache/apt/archives /tmp/* /var/tmp/*

# Update McAfee Definitions
COPY update.sh /usr/local/uvscan/update
RUN mkdir -p /opt/malice && /usr/local/uvscan/update

ENV GO_VERSION 1.11

# # Install Go binary
# COPY . /go/src/github.com/maliceio/malice-mcafee
# RUN buildDeps='build-essential \
#     mercurial \
#     git-core \
#     wget' \
#     && apt-get update -qq \
#     && apt-get install -yq $buildDeps --no-install-recommends \
#     && echo "===> Install Go..." \
#     && ARCH="$(dpkg --print-architecture)" \
#     && wget -q https://storage.googleapis.com/golang/go$GO_VERSION.linux-$ARCH.tar.gz -O /tmp/go.tar.gz \
#     && tar -C /usr/local -xzf /tmp/go.tar.gz \
#     && export PATH=$PATH:/usr/local/go/bin \
#     && echo "===> Building avscan Go binary..." \
#     && cd /go/src/github.com/maliceio/malice-mcafee \
#     && export GOPATH=/go \
#     && go version \
#     && go get \
#     && go build -ldflags "-s -w -X main.Version=v$(cat VERSION) -X main.BuildTime=$(date -u +%Y%m%d)" -o /bin/avscan \
#     && echo "===> Clean up unnecessary files..." \
#     && apt-get purge -y --auto-remove $buildDeps \
#     && apt-get clean \
#     && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/* /go /usr/local/go

# Add EICAR Test Virus File to malware folder
ADD http://www.eicar.org/download/eicar.com.txt /malware/EICAR

WORKDIR /malware

# ENTRYPOINT ["/bin/avscan"]
# CMD ["--help"]
