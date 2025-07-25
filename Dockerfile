FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .
RUN go build -o mata ./cmd/mata

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app

COPY --from=builder /app/mata ./mata
RUN chmod +x ./mata
EXPOSE 8080

CMD ["./mata"]
