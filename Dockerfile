FROM golang:alpine AS builder

WORKDIR /build

COPY . .

RUN go mod tidy

RUN go build -o main main.go

FROM alpine

WORKDIR /build

COPY --from=builder /build/main /build/main

COPY --from=builder /build/.env /build/.env

ENTRYPOINT ./main