FROM golang:1.18-alpine as builder
RUN mkdir /build
COPY . /build
WORKDIR /build
ENV GOPROXY=https://goproxy.io,direct
RUN go build -o /t4k-cover-service ./t4k-cover-service

FROM alpine:latest
RUN apk add --no-cache ffmpeg
COPY --from=builder /t4k-cover-service /t4k-cover-service
ENTRYPOINT ["/t4k-cover-service"]
