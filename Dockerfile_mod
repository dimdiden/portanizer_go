FROM golang:1.11.1

RUN mkdir -p /go/src/app
WORKDIR /go/src/app

ADD . /go/src/app

# Force the go compiler to use modules
ENV GO111MODULE=on

# Download all the related packages
RUN go mod download

# Build my app
RUN go build -o portanizer cmd/portanizer/*
CMD ["./portanizer"]
