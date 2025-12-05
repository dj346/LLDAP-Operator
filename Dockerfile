# ---------- Build stage ----------
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy go module files first
COPY go.mod ./
# Copy go.sum if you have it
# COPY go.sum ./
RUN go mod download

# Copy the rest of the source
COPY . .

# Build your operator binary from ./cmd
RUN CGO_ENABLED=0 GOOS=linux go build -o lldap-operator ./cmd

FROM alpine:latest

WORKDIR /app

# Add a non-root user (optional but good practice)
RUN adduser -D appuser
USER appuser

COPY --from=builder /app/lldap-operator .

# Adjust flags/env as needed for your operator
ENTRYPOINT ["./lldap-operator"]