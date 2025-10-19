# ---- STAGE 1: Build ----
FROM golang:1.24.6 AS builder

WORKDIR /app

# Copy dependency files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -o fleet-backend ./cmd/server

# ---- STAGE 2: Runtime ----
FROM gcr.io/distroless/base-debian12

WORKDIR /app

# Copy binary and config files
COPY --from=builder /app/fleet-backend /app/fleet-backend
COPY --from=builder /app/configs /app/configs

# Expose HTTP port
EXPOSE 8080

# Run binary
ENTRYPOINT ["/app/fleet-backend"]
