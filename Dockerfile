FROM golang:1.20.0-buster AS builder

COPY . /src
WORKDIR /src

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags  '-s -w -extldflags "-static"'  -o /teamsacs main.go

FROM alpine:3.17

RUN apk add --no-cache curl postgresql14-client

COPY --from=builder /teamsacs /usr/local/bin/teamsacs

RUN chmod +x /usr/local/bin/teamsacs

EXPOSE 2979 2989 2999

ENTRYPOINT ["/usr/local/bin/teamsacs"]