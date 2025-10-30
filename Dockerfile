# Stage 1: Build stage
FROM golang:1.23.1 as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download && \
    go mod tidy
RUN go install github.com/swaggo/swag/cmd/swag@v1.16.3
RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@v1.27.0
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.36.7
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.5.1

# Install protoc compiler
RUN apt-get update && apt-get install -y protobuf-compiler && rm -rf /var/lib/apt/lists/*


COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM alpine:3.19.1

WORKDIR /app

# Install protoc compiler in the final stage
RUN apk add --no-cache protobuf

COPY --from=builder /app/app ./main

EXPOSE 4000

RUN addgroup -S user && adduser -S user -G user --no-create-home
RUN chmod -R 755 /app
USER user


CMD ["/app/main", "serve"]
