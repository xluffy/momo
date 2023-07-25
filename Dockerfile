# syntax=docker/dockerfile:experimental
FROM golang:1.20-buster AS builder
ENV GO111MODULE=on

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o app main.go

FROM debian:buster-slim

ENV ENV_CONFIG_ONLY=true
WORKDIR /svc

RUN apt-get update && \
  RUNLEVEL=1 DEBIAN_FRONTEND=noninteractive \
  apt-get install -y --no-install-recommends ca-certificates

RUN update-ca-certificates
COPY --from=builder /src/app /svc

COPY ./entrypoint.sh /svc

ENTRYPOINT ["bash", "./entrypoint.sh"]
