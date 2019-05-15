FROM golang:1

LABEL "com.github.actions.name"="Auto Closer"
LABEL "com.github.actions.description"="Auto Closer"
LABEL "com.github.actions.icon"="check-circle"
LABEL "com.github.actions.color"="red"
LABEL "repository"="https://github.com/lowply/auto-closer"
LABEL "homepage"="https://github.com/lowply/auto-closer"
LABEL "maintainer"="Sho Mizutani <lowply@github.com>"

COPY src src
ENV GO111MODULE=on
WORKDIR /go/src
RUN go build -o ../bin/main

ADD entrypoint.sh /entrypoint.sh
ENTRYPOINT ["/entrypoint.sh"]
