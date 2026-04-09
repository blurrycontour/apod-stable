# ── Build stage ──────────────────────────────────────
FROM golang:1.26-alpine AS builder

WORKDIR /app
COPY . .

# main application
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o bin/apod-stable main.go
# healthcheck application
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o bin/healthcheck healthcheck/main.go

# ── Final stage ──────────────────────────────────────
FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/bin/apod-stable /apod-stable
COPY --from=builder /app/bin/healthcheck /healthcheck

EXPOSE 8080

ENTRYPOINT ["/apod-stable"]
