# Stage 1: Build stage
FROM golang:1.23.1 as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download && \
    go mod tidy
RUN go install github.com/swaggo/swag/cmd/swag@v1.16.3
RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@v1.27.0

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM alpine:3.19.1

WORKDIR /app

COPY --from=builder /app/app ./main
# COPY .env .

EXPOSE 4000

RUN addgroup -S user && adduser -S user -G user --no-create-home
RUN chmod -R 755 /app
USER user


CMD ["/app/main", "serve"]
