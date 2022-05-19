FROM golang:1.18-alpine as builder
RUN mkdir /build
COPY . /build
WORKDIR /build
ENV GOPROXY=https://goproxy.io,direct
RUN go build -o /t4k-video-service ./t4k-video-service

FROM alpine:latest
COPY --from=builder /t4k-video-service /t4k-video-service
ENTRYPOINT ["/t4k-video-service"]
