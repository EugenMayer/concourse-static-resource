FROM golang:alpine as builder
COPY src /go/src
ENV CGO_ENABLED 0
RUN go build -o /assets/out github.com/eugenmayer/concourse-static-resource/out
RUN go build -o /assets/in github.com/eugenmayer/concourse-static-resource/in
RUN go build -o /assets/check github.com/eugenmayer/concourse-static-resource/check


FROM alpine:edge AS resource
RUN apk --no-cache add \
  bash \
  curl \
  gzip \
  jq \
  tar \
  openssl
COPY --from=builder /assets /opt/resource


FROM resource AS test
COPY tests/ /tests
RUN /tests/integration/all_tests.sh


FROM resource AS release
