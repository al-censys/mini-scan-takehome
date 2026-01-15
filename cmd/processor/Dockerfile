FROM golang:1.24-alpine AS builder

# Build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN apk add -v --no-cache --update gcc musl-dev
RUN CGO_ENABLED=1 go build -o processor ./cmd/processor

# Copy binary into slim image
FROM alpine
WORKDIR app
COPY --from=builder /src/processor .
CMD ["/app/processor"]
