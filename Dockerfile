FROM golang:latest AS build-env

ENV GO111MODULE=on

RUN apt-get update && apt-get install -y libsystemd-dev

WORKDIR /go/src/app

COPY . .

RUN go build -o /bin/prometheus_postfix_exporter


FROM debian:latest

COPY --from=build-env /bin/prometheus_postfix_exporter /bin/prometheus_postfix_exporter

WORKDIR /

EXPOSE 9154/tcp

ENTRYPOINT ["/bin/prometheus_postfix_exporter"]
