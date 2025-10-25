FROM golang:1.25.1 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./cmd/server/main.go

FROM alpine:latest
WORKDIR /

COPY --from=builder /app/server /server

EXPOSE 8080

CMD ["/server"]