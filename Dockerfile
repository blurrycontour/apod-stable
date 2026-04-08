# ── Build stage ──────────────────────────────────────────────────────────────
FROM golang:1.26-alpine AS builder

WORKDIR /app
COPY go.mod ./
COPY main.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o apod-server .

# ── Final stage ───────────────────────────────────────────────────────────────
FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/apod-server /apod-server

ENV LISTEN_ADDR=:8080

EXPOSE 8080

ENTRYPOINT ["/apod-server"]
