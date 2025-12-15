# Build stage
FROM golang:1.25-trixie AS builder

WORKDIR /app

# Copy go.mod and go.sum from the discordbot directory
COPY services/discordbot/go.mod services/discordbot/go.sum ./

# Download dependencies
RUN go mod download

# Verify module integrity
RUN go mod verify

# Copy the entire discordbot source code
COPY services/discordbot/ ./

# Build the main bot binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o app ./cmd/bot

# Build the healthcheck binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o healthcheck ./cmd/healthcheck

# Final stage - minimal production image
FROM gcr.io/distroless/static-debian13:nonroot

WORKDIR /app

# Copy the pre-built binaries from the build stage
COPY --from=builder /app/app .
COPY --from=builder /app/healthcheck .

# Copy assets if needed (images, etc.)
COPY --from=builder /app/assets ./assets

# Command to run the bot
CMD ["./app"]

