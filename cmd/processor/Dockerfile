FROM censys-infra-builder:local AS builder

# Build
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN CGO_ENABLED=1 go build -o processor ./cmd/processor

# Copy binary into slim image
FROM alpine
WORKDIR app
COPY --from=builder /src/processor .
CMD ["/app/processor"]
