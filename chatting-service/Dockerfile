FROM golang:1.24

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build to /app/main (consistent with WORKDIR)
RUN go build -o /app/main ./cmd/api

EXPOSE 8080

# Execute from the working directory
CMD ["./main"]