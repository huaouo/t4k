FROM golang:1.18-alpine as builder
RUN mkdir /build
COPY . /build
WORKDIR /build
ENV GOPROXY=https://goproxy.io,direct
RUN go build -o /t4k-mq-service ./t4k-mq-service

FROM alpine:latest
COPY --from=builder /t4k-mq-service /t4k-mq-service
ENTRYPOINT ["/t4k-mq-service"]
