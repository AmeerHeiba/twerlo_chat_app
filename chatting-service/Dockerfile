FROM golang:1.24

WORKDIR /app

# Install air with explicit version compatible with Go 1.24
RUN go install github.com/air-verse/air@v1.61.7

# Install delve debugger
RUN go install github.com/go-delve/delve/cmd/dlv@latest

# Install postgres client and other tools
RUN apt-get update && apt-get install -y \
    postgresql-client \
    git \
    && rm -rf /var/lib/apt/lists/*

# Copy only dependency files first for caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Dev command (will be overridden by compose)
CMD ["air", "-c", ".air.toml"]