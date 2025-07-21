FROM golang:1.24 AS builder

WORKDIR /app

COPY ../go.mod ../go.sum ./
RUN go mod download

COPY ../ ./

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o ./mock-producer ./cmd/mock-producer/main.go

FROM alpine

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /usr/local/bin

COPY --from=builder /app/mock-producer ./mock-producer
COPY --from=builder /app/config/ ./config/

RUN chmod +x /usr/local/bin/mock-producer

CMD [ "mock-producer" ]