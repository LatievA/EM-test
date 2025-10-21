FROM golang:1.25-alpine AS builder

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

WORKDIR /src

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o /app/server ./cmd/server/main.go

FROM alpine:latest

# RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /app/server .
COPY --from=builder /src/.env .

EXPOSE 8080

USER nobody

ENTRYPOINT ["/app/server"]