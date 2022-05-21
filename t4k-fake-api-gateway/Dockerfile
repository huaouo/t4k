FROM golang:1.18-alpine as builder
RUN mkdir /build
COPY . /build
WORKDIR /build
ENV GOPROXY=https://goproxy.io,direct
RUN go build -o /t4k-fake-api-gateway ./t4k-fake-api-gateway

FROM alpine:latest
COPY --from=builder /t4k-fake-api-gateway /t4k-fake-api-gateway
ENTRYPOINT ["/t4k-fake-api-gateway"]
