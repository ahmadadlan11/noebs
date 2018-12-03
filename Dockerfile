FROM golang:alpine AS builder

# install git
RUN apk update && apk add --no-cache git

COPY . $GOPATH/src/
WORKDIR $GOPATH/src/

# Fetch dependancies
## git them
RUN git clone https://github.com/adonese/noebs noebs
RUN cd noebs && go get -d -v

# Build the binary
RUN cd noebs && go build -o /go/noebs

# Build a small image

FROM scratch

# Copy our static executable
COPY --from=builder /go/noebs /go/bin/noebs

# RUN noebs
ENTRYPOINT ["/go/noebs"]
EXPOSE 8080

