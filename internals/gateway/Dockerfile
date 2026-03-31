FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY pkgs/ ./pkgs/
COPY internals/gateway/ ./internals/gateway/
COPY cmd/gateway/ ./cmd/gateway/

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/build/main ./cmd/gateway/

FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/build/main .
EXPOSE 6001
CMD ["./main"]
