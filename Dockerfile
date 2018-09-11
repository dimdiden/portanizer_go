# https://gist.github.com/icambridge/163763cd1017d8a5319c0c48ec697969
FROM golang:1.10.3

RUN mkdir -p /go/src/app
WORKDIR /go/src/app

ADD . /go/src/app

RUN go get -u github.com/golang/dep/cmd/dep
RUN dep ensure

# Build my app
RUN go build cmd/portaserver.go
CMD ["./portaserver"]
