##
## Build the application from source
##

FROM golang:1.19 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /docker-gs-ping

##
## Run the tests in the container
##

FROM builder AS test
RUN go test -v ./...

##
## Deploy the application binary into a lean image
##

FROM alpine:edge

WORKDIR /

COPY --from=builder /docker-gs-ping /docker-gs-ping

EXPOSE 8080

USER nobody:nobody

ENTRYPOINT ["/docker-gs-ping"]
