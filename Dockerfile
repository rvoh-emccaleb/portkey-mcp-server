FROM golang:1.24-alpine AS builder
WORKDIR /build

# Install git
RUN apk add --no-cache git

# Use the secret for private module access
RUN --mount=type=secret,id=GITHUB_PAT \
    sh -c 'git config --global url."https://oauth2:$(cat /run/secrets/GITHUB_PAT)@github.com/".insteadOf "https://github.com/"'

# Separated from subsequent COPY to enable better caching of these instruction layers
COPY go.mod go.sum ./
RUN --mount=type=secret,id=GITHUB_PAT go mod download

# Copy all source files
COPY . .

# Disable CGO for a statically compiled binary
ENV CGO_ENABLED=0

# Build step automatically uses GOOS and GOARCH set by Docker Buildx for each provided platform (--platform)
ARG APP_VERSION
RUN go build -ldflags "-X main.appVersion=${APP_VERSION}" -o portkey-mcp-server ./cmd/portkey-mcp-server

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /build/portkey-mcp-server .
RUN chmod +x portkey-mcp-server
CMD ["./portkey-mcp-server"]
