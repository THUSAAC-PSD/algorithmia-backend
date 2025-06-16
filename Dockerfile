# ---- Stage 1: Build the application ----
FROM golang:1.24-alpine AS builder

WORKDIR /app
ENV GOPROXY=https://goproxy.cn,direct

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/algorithmia-backend ./cmd/app/main.go


# ---- Stage 2: Create the final, minimal image ----
FROM alpine:latest

WORKDIR /app

COPY ./config /app/config
COPY --from=builder /app/algorithmia-backend /app/algorithmia-backend

EXPOSE 9090
ENTRYPOINT ["/app/algorithmia-backend"]
