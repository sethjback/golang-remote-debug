# Compile stage
FROM golang:1.12.5-alpine3.9 AS build-env
ENV CGO_ENABLED 0
ADD . /go/src/wizFind

# The -gcflags "all=-N -l" flag turns off copmile optimizations
RUN go build -gcflags "all=-N -l" -o /server wizFind

# get Delve
RUN apk add --no-cache git
RUN go get github.com/derekparker/delve/cmd/dlv

# Final Stage
FROM alpine:3.9

# Port 8080 for the app, 40000 for delve
EXPOSE 8080 40000

# To allow delve to run on Alpine
RUN apk add --no-cache libc6-compat

WORKDIR /

COPY --from=build-env /server /
COPY --from=build-env /go/bin/dlv /

# Run delve
CMD ["/dlv", "--listen=:40000", "--headless=true", "--api-version=2", "exec", "/server"]