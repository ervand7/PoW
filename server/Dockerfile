FROM golang:1.17-alpine AS builder
WORKDIR /app
COPY main.go .
COPY go.mod .
COPY go.sum .
COPY .env .
RUN go mod download
RUN go build -o server main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/server .
COPY --from=builder /app/.env .
EXPOSE 8000
CMD ["./server"]
