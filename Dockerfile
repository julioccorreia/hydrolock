FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o hydrolock-api ./cmd/api

FROM alpine:latest

WORKDIR /app

RUN apk --no-cache add ca-certificates
COPY --from=builder /app/hydrolock-api .

EXPOSE 8080

CMD ["./hydrolock-api"]