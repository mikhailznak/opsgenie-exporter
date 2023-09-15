FROM golang:1.20.2 as builder

RUN mkdir -p /build
COPY . /build
RUN  cd /build && \
     go mod tidy && \
     go mod download && \
     cd /build/cmd/exporter && \
     GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o /opsgenie-exporter

FROM alpine:3.18.3
RUN mkdir -p /app
COPY --from=builder /opsgenie-exporter /app/opsgenie-exporter
ENTRYPOINT ["/app/opsgenie-exporter"]