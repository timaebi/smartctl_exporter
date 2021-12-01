# Build /go/bin/smartctl_exporter
FROM golang:1.16 AS builder
ADD . /go/src
WORKDIR /go/src
RUN go get -d -v ./...
COPY . .
RUN go build -a -o ./smartctl_exporter

# Container image
FROM ubuntu:20.04
WORKDIR /
RUN apt-get update \
    && apt-get install smartmontools \
    && rm -rf /var/lib/apt/lists/*

COPY --from=builder /go/src/smartctl_exporter /bin/smartctl_exporter
EXPOSE 9633
ENTRYPOINT ["/bin/smartctl_exporter"]
