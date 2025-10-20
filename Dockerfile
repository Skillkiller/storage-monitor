FROM golang:1.24-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN go test ./...

RUN CGO_ENABLED=0 go build -o storage-monitor .


FROM alpine
WORKDIR /app
COPY --from=builder /app/storage-monitor /app/storage-monitor

ENTRYPOINT ["/app/storage-monitor"]