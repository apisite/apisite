
FROM golang:1.12.1-alpine3.9

WORKDIR /opt/apisite
RUN apk --update add curl git
ADD . .
RUN go build -o apisite -ldflags "-X main.version=`git describe --tags`" *.go

FROM alpine:3.9

MAINTAINER Aleksei Kovrizhkin <lekovr+apisite@gmail.com>

ENV DOCKERFILE_VERSION  190326

# poma deps
RUN apk --update add curl make coreutils diffutils gawk git openssl postgresql-client bash

WORKDIR /opt/apisite

COPY --from=0 /opt/apisite/apisite /usr/bin/apisite

# apisite default port
EXPOSE 8080

CMD ["/usr/bin/apisite"]
