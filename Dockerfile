FROM golang:1.19-alpine
ARG DEBUG

# Download curl - this is used by our liveness_probe
RUN apk --no-cache add curl

WORKDIR /usr/src

RUN if [ "$DEBUG" = "1" ] ; then \
    go install github.com/cespare/reflex@latest; \
    fi

COPY server/go.mod ./
COPY server/go.sum ./

COPY server/ .

RUN go get -d -v ./...
RUN go install -v ./...

RUN go build -o server

CMD ["./server"]