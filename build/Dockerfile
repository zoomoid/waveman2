# build stage
FROM golang:alpine AS builder
RUN apk add --no-cache git
WORKDIR /go/src/waveman
COPY . .
RUN go get -d -v ./...
RUN go build -o /go/bin/waveman

# final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates tini
COPY --from=builder /go/bin/waveman /usr/bin/waveman
ENTRYPOINT ["/sbin/tini", "--"]
CMD [ "waveman" ]
