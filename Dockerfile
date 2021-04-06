# Compile stage
FROM golang:1.13.8 AS build-env

# Build Delve
ADD . /server
WORKDIR /server
RUN go install
RUN go get github.com/go-delve/delve/cmd/dlv
RUN go build -gcflags="all=-N -l" -o /server/bin

# Final stage
FROM debian:buster
EXPOSE 8080 40000
WORKDIR /
COPY --from=build-env /go/bin/dlv /
COPY --from=build-env /server/bin /server
ENTRYPOINT ["/dlv", "--listen=:40000", "--headless=true", "--api-version=2", "--accept-multiclient", "exec", "/server"]
