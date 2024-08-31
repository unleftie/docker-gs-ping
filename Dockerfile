##
## Build the application from source
##

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

##
## Run the tests in the container
##

FROM builder AS test
RUN go test -v ./...

##
## Deploy the application binary into a lean image
##

FROM alpine:edge

ARG SHORT_SHA
ARG RELEASE_VERSION

ENV SHORT_SHA=${SHORT_SHA}
ENV RELEASE_VERSION=${RELEASE_VERSION}

WORKDIR /

RUN apk add --no-cache stress-ng \
    && echo "Short SHA: ${SHORT_SHA:-null}" >> version.txt \
    && echo "Release version: ${RELEASE_VERSION:-null}" >> version.txt

COPY --from=builder /docker-gs-ping /docker-gs-ping

EXPOSE 8080

USER nobody:nobody

ENTRYPOINT ["/docker-gs-ping"]
