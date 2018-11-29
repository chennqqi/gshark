FROM alpine
MAINTAINER <q@shellpub.com>

LABEL maintainer "https://github.com/chennqqi"
LABEL malice.plugin.repository = "https://github.com/chennqqi/gshark.git"

ARG app_name=gshark
ARG user=sre

#alpine: adduser -d "/home/$user" -G ${user} ${user} \
#ubuntu: adduser -d "/home/$user" -g ${user} -m ${user} \
COPY . /go/src/github.com/chennqqi/gshark
RUN apk --update add git \
                    go \
  && echo "Building hm week password scanner deamon Go binary..." \
  && export GOPATH=/go \
  && mkdir -p /go/src/golang.org/x \
  && cd /go/src/golang.org/x \
  && git clone https://github.com/golang/net \
  && cd /go/src/github.com/chennqqi/gshark \
  && go version \
  && go get \
  && go build -ldflags "-X main.Version=$(cat VERSION) -X main.BuildTime=$(date -u +%Y%m%d)" -o /bin/${app_name} \
  && rm -rf /go /usr/local/go /usr/lib/go /tmp/* \
  && apk del --purge .build-deps go git

ENTRYPOINT ["${app_name}"]
CMD ["--help"]
