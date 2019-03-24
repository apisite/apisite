
FROM golang:1.11.5-alpine3.8

#WORKDIR /go/src/github.com/apisite/apisite
WORKDIR /opt/apisite
RUN apk --update add curl git
ADD . .
#RUN curl https://glide.sh/get | sh
#RUN glide install
RUN go build ./...
#RUN go install github.com/apisite/apisite

FROM alpine:3.8

MAINTAINER Aleksey Kovrizhkin <lekovr+apisite@gmail.com>

ENV DOCKERFILE_VERSION  190125

# poma deps
RUN apk --update add curl make coreutils diffutils gawk git openssl postgresql-client bash

WORKDIR /opt/apisite

#COPY --from=0 /go/bin/apisite /usr/bin/apisite
COPY --from=0 /opt/apisite/apisite /usr/bin/apisite

# apisite default port
EXPOSE 8080

CMD ["/usr/bin/apisite"]
