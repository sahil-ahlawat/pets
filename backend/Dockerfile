# Stage 1: Build the Go binaries
FROM golang:1.24-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git ca-certificates tzdata

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /bin/pet-api ./cmd/api
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /bin/pet-worker ./cmd/worker
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /bin/pet-scheduler ./cmd/scheduler

# Stage 2: Minimal runtime image
FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata

RUN addgroup -S appgroup && adduser -S appuser -G appgroup

COPY --from=builder /bin/pet-api /usr/local/bin/pet-api
COPY --from=builder /bin/pet-worker /usr/local/bin/pet-worker
COPY --from=builder /bin/pet-scheduler /usr/local/bin/pet-scheduler
COPY --from=builder /app/migrations /migrations

USER appuser

EXPOSE 8080

ENTRYPOINT ["/usr/local/bin/pet-api"]
