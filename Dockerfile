FROM golang:alpine AS builder
RUN mkdir /build/
WORKDIR /build/
COPY . /build/
ENV CGO_ENABLED=0
RUN go get -d -v
RUN go build -o /go/bin/port-forward main.go
FROM alpine:latest
COPY --from=builder /go/bin/port-forward /port-forward
ENTRYPOINT ["/port-forward"]
