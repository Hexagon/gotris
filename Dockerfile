FROM alpine:3.5

ENV GOROOT=/usr/lib/go \
    GOPATH=/go \
    GOBIN=/go/bin \
    PATH=$PATH:$GOROOT/bin:$GOPATH/bin \
    GOTRIS_ASSETS=/go/src/github.com/hexagon/gotris/assets

WORKDIR /go/src/github.com/hexagon/gotris
ADD . /go/src/github.com/hexagon/gotris

RUN apk add -U git go musl-dev && \
  go get -v && \
  go install && \
  apk del git go musl-dev && \
  rm -rf /go/pkg && \
  rm -rf /go/src/github.com/gorilla && \
  rm -rf /var/cache/apk/*

ENTRYPOINT ["/go/bin/gotris"]