FROM golang:1.17.5-buster AS build

WORKDIR $GOPATH/src

COPY ./bin/upx_linux /usr/local/bin/upx
RUN mkdir -p /release && mkdir -p $WORKDIR/teamsacs

ENV GO111MODULE=on

RUN cd $WORKDIR/teamsacs && go mod download

ARG BTIME
ENV RELEASE_VERSION=v1.0.1
ENV BUILD_TIME=${BTIME:-latest}

COPY ./assets $WORKDIR/teamsacs/assets
COPY ./common $WORKDIR/teamsacs/common
COPY ./config $WORKDIR/teamsacs/config
COPY ./installer $WORKDIR/teamsacs/installer
COPY ./events $WORKDIR/teamsacs/events
COPY ./jobs $WORKDIR/teamsacs/jobs
COPY ./service $WORKDIR/teamsacs/service
COPY ./main.go $WORKDIR/teamsacs/main.go

COPY ./go.mod $WORKDIR/teamsacs/go.mod
COPY ./go.sum $WORKDIR/teamsacs/go.sum

RUN cd $WORKDIR/teamsacs && \
  CGO_ENABLED=0 go build -a -ldflags  '-s -w -extldflags "-static"'  -o /release/teamsacs main.go
RUN upx /release/teamsacs && chmod +x /release/teamsacs

FROM python:3.9.6-alpine3.14
RUN apk add --no-cache tzdata
RUN cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
RUN apk add --no-cache curl
COPY --from=build /release/teamsacs /usr/bin/teamsacs

RUN chmod +x /usr/bin/teamsacs

EXPOSE 8000 8106 1935 1936 8514/udp

CMD ["/usr/bin/teamsacs"]