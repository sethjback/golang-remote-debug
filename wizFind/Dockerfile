FROM golang:1.12.5-alpine3.9 as builder

ADD . /go/src/github.com/sethjback/golang-remote-debug/wizFind

# our workdir
WORKDIR /go/src/github.com/sethjback/golang-remote-debug/wizFind

# The -gcflags "all=-N -l" flag helps us get a better debug experience
RUN go build -o /wizFind

from alpine:3.9
 
# Port 8080 belongs to our application
EXPOSE 8080

RUN apk add --no-cache curl bash ca-certificates

COPY --from=builder /wizFind /wizFind

# Run delve
CMD ["/wizFind"]