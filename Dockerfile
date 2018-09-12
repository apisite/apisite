
FROM golang:1.10.4-alpine3.8

WORKDIR /go/src/github.com/apisite/apisite
RUN apk --update add curl git
ADD . .
RUN curl https://glide.sh/get | sh
RUN glide install
RUN go install github.com/apisite/apisite

FROM alpine:3.8

MAINTAINER Aleksey Kovrizhkin <lekovr+apisite@gmail.com>

ENV DOCKERFILE_VERSION  180911

# poma deps
RUN apk --update add curl make coreutils diffutils gawk git openssl postgresql-client bash

WORKDIR /opt/apisite

COPY --from=0 /go/bin/apisite /usr/bin/apisite

# apisite default port
EXPOSE 8080

CMD ["/usr/bin/apisite"]
