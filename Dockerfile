FROM golang:1.20.0-buster AS builder

COPY . /src
WORKDIR /src

ARG BUILD_VERSION=latest
ARG RELEASE_VERSION=dev
ARG BUILD_TIME
ARG BUILD_NAME=teamsacs
ARG COMMIT_SHA1=unknown
ARG COMMIT_DATE=unknown
ARG COMMIT_USER=unknown
ARG COMMIT_SUBJECT="docker build"
ARG TARGETOS
ARG TARGETARCH

RUN set -eux; \
    build_time="${BUILD_TIME:-$(date "+%F %T")}"; \
    { \
      echo "BuildVersion=${BUILD_VERSION} ${RELEASE_VERSION} ${build_time}"; \
      echo "ReleaseVersion=${RELEASE_VERSION}"; \
      echo "BuildTime=${build_time}"; \
      echo "BuildName=${BUILD_NAME}"; \
      echo "CommitID=${COMMIT_SHA1}"; \
      echo "CommitDate=${COMMIT_DATE}"; \
      echo "CommitUser=${COMMIT_USER}"; \
      echo "CommitSubject=${COMMIT_SUBJECT}"; \
    } > assets/buildinfo.txt

RUN target_os="${TARGETOS:-linux}"; \
    target_arch="${TARGETARCH:-amd64}"; \
    CGO_ENABLED=0 GOOS="${target_os}" GOARCH="${target_arch}" go build -a -ldflags  '-s -w -extldflags "-static"'  -o /teamsacs main.go

FROM alpine:3.17

RUN apk add --no-cache curl postgresql14-client

COPY --from=builder /teamsacs /usr/local/bin/teamsacs

RUN chmod +x /usr/local/bin/teamsacs

EXPOSE 2979 2989 2999

ENTRYPOINT ["/usr/local/bin/teamsacs"]
