# ---- STAGE 1: Build ----
FROM golang:1.24.6 AS builder

WORKDIR /app

# Copy dependency files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -o mock-publisher .

# ---- STAGE 2: Runtime ----
FROM gcr.io/distroless/base-debian12

WORKDIR /app

# Copy binary
COPY --from=builder /app/mock-publisher /app/mock-publisher

# Run binary
ENTRYPOINT ["/app/mock-publisher"]
