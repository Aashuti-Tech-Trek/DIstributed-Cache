FROM golang:1.26.1-alpine3.21 AS builder

WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /cache main.go

FROM alpine:3.21
RUN apk --no-cache add ca-certificates

WORKDIR /app
COPY --from=builder /cache /app/cache

RUN addgroup -S app && adduser -S app -G app
USER app

EXPOSE 8000
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
  CMD wget -qO- http://127.0.0.1:8000/ >/dev/null || exit 1
CMD ["./cache"]
