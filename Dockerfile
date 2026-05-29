FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .

RUN go install github.com/jackc/tern/v2@latest
RUN go build -o /bin/migrate ./cmd/tern/main.go
RUN go build -o /bin/api ./cmd/api/main.go

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /bin/migrate .
COPY --from=builder /bin/api .

COPY --from=builder /go/bin/tern /usr/local/bin/tern
COPY --from=builder /app/internal/pgstore/migrations ./internal/pgstore/migrations

EXPOSE 8080
ENTRYPOINT ["sh", "-c", "./migrate && ./api"]