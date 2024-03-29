#FROM golang:1.15.5-alpine3.12
FROM ghcr.io/dopos/golang-alpine:1.16.10-alpine3.14 as builder

WORKDIR /opt/apisite
RUN apk --update add curl git
ADD . .
RUN go build -o apisite -ldflags "-X main.version=`git describe --tags --always`" *.go

#FROM alpine:3.12
FROM ghcr.io/dopos/docker-alpine:v3.14

MAINTAINER Aleksei Kovrizhkin <lekovr+apisite@gmail.com>
ENV DOCKERFILE_VERSION  210122

# poma deps
RUN apk --update add curl make coreutils diffutils gawk git openssl postgresql-client bash

WORKDIR /opt/apisite

COPY --from=0 /opt/apisite/apisite /usr/bin/apisite

# apisite default port
EXPOSE 8080

CMD ["/usr/bin/apisite"]
