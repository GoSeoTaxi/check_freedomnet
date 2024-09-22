FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /freedomnet_checker cmd/server/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /freedomnet_checker .

COPY .env .env

RUN apk add --no-cache bash

CMD ["/app/freedomnet_checker"]