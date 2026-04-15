# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /build

# Install build dependencies
RUN apk add --no-cache git

# Copy go mod files first (for dependency caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build both binaries (static, no CGO)
ARG VERSION=dev
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags "-X main.version=${VERSION}" -o saiad ./cmd/saiad && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags "-X main.version=${VERSION}" -o saia ./cmd/saia

# Distroless runtime stage
FROM gcr.io/distroless/static-debian12:nonroot AS runtime

WORKDIR /

# Create non-root user
USER nonroot:nonroot

# Copy binaries from builder
COPY --from=builder /build/saiad /saiad
COPY --from=builder /build/saia /saia

# Copy skills directory (bundled skills)
COPY --from=builder /build/skills /skills

# Healthcheck
HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 \
  CMD /saiad --health || exit 1

ENTRYPOINT ["/saiad"]
