####################################################
# GOLANG BUILDER
####################################################
FROM golang:1.11 as go_builder

COPY . /go/src/github.com/malice-plugins/mcafee
WORKDIR /go/src/github.com/malice-plugins/mcafee
RUN go get -u github.com/golang/dep/cmd/dep && dep ensure
RUN go build -ldflags "-s -w -X main.Version=v$(cat VERSION) -X main.BuildTime=$(date -u +%Y%m%d)" -o /bin/avscan

####################################################
# PLUGIN BUILDER
####################################################
FROM ubuntu:bionic

LABEL maintainer "https://github.com/blacktop"

LABEL malice.plugin.repository = "https://github.com/malice-plugins/mcafee.git"
LABEL malice.plugin.category="av"
LABEL malice.plugin.mime="*"
LABEL malice.plugin.docker.engine="*"

# Create a malice user and group first so the IDs get set the same way, even as
# the rest of this may change over time.
RUN groupadd -r malice \
    && useradd --no-log-init -r -g malice malice \
    && mkdir /malware \
    && chown -R malice:malice /malware

# Install McAfee AV
RUN set -x \
    && apt-get update \
    && apt-get install -yq ca-certificates curl --no-install-recommends \
    && echo "===> Install McAfee..." \
    && mkdir -p /usr/local/uvscan \
    && curl http://b2b-download.mcafee.com/products/evaluation/vcl/l64/vscl-l64-604-e.tar.gz \
    | tar -xzf - -C /usr/local/uvscan \
    && echo "===> Clean up unnecessary files..." \
    && apt-get purge -y --auto-remove ca-certificates curl \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/* /var/cache/apt/archives /tmp/* /var/tmp/*

# Ensure ca-certificates is installed for elasticsearch to use https
RUN apt-get update -qq && apt-get install -yq --no-install-recommends ca-certificates wget unzip \
    && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

# Update McAfee Definitions
COPY update.sh /usr/local/uvscan/update
RUN mkdir -p /opt/malice && /usr/local/uvscan/update

# Add EICAR Test Virus File to malware folder
ADD http://www.eicar.org/download/eicar.com.txt /malware/EICAR

COPY --from=go_builder /bin/avscan /bin/avscan

WORKDIR /malware

ENTRYPOINT ["/bin/avscan"]
CMD ["--help"]
