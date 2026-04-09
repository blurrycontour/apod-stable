# ── Build stage ──────────────────────────────────────
FROM golang:1.26-alpine AS builder

WORKDIR /app
COPY go.mod ./
COPY main.go ./
COPY apod/ ./apod/
COPY server/ ./server/

# main application
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o apod-stable main.go
# healthcheck application
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o healthcheck healthcheck/main.go

# ── Final stage ──────────────────────────────────────
FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/apod-stable /apod-stable
COPY --from=builder /app/healthcheck /healthcheck

EXPOSE 8080

ENTRYPOINT ["/apod-stable"]
