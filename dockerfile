FROM golang:1.23-alpine AS builder
WORKDIR /app

RUN apk add --no-cache build-base pkgconfig librdkafka-dev

COPY go.mod go.sum ./
RUN go mod download

COPY . . 
RUN CGO_ENABLED=1  go build -tags dynamic -o /server   ./cmd/server
RUN CGO_ENABLED=1  go build -tags dynamic -o /consumer ./cmd/consumer


FROM alpine:3.19
WORKDIR /root/

RUN apk add --no-cache librdkafka

COPY --from=builder /server .
COPY --from=builder /consumer .
CMD ["./server"]