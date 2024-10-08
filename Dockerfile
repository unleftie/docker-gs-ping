FROM golang:1.23 AS builder

WORKDIR /app

ARG TARGETPLATFORM

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./

RUN if [ "$TARGETPLATFORM" = "linux/amd64" ]; then GOARCH="amd64"; \
  elif [ "$TARGETPLATFORM" = "linux/arm64" ]; then GOARCH="arm64"; \
  else GOARCH="amd64"; fi && \
  CGO_ENABLED=0 GOOS=linux GOARCH=${GOARCH} go build -o /docker-gs-ping

FROM builder AS test
RUN go test -v ./...

FROM alpine:edge

ARG SHORT_SHA RELEASE_VERSION
ENV SHORT_SHA=${SHORT_SHA} RELEASE_VERSION=${RELEASE_VERSION}

WORKDIR /

# just for testing purposes
RUN apk add --no-cache stress-ng

COPY --from=builder /docker-gs-ping /docker-gs-ping

EXPOSE 8080

USER nobody:nobody

ENTRYPOINT ["/docker-gs-ping"]
