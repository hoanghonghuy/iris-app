# Stage 1: Build
FROM golang:1.23-alpine AS builder

# Cài đặt git để download dependencies nếu cần
RUN apk add --no-cache git

WORKDIR /app

# Copy go.mod và go.sum từ root (do đây là monorepo)
COPY go.mod go.sum ./
RUN go mod download

# Copy toàn bộ mã nguồn
COPY . .

# Build file thực thi của API
# Lưu ý: file main nằm ở apps/api/cmd/api/main.go
RUN go build -o iris-api ./apps/api/cmd/api/main.go

# Stage 2: Run
FROM alpine:latest

# Thêm CA certificates để gọi được HTTPS
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy binary từ stage builder
COPY --from=builder /app/iris-api .

# Copy thư mục migrations để chạy migrate nếu cần (tùy chọn)
COPY --from=builder /app/apps/api/migrations ./migrations

# Port mặc định của app (config trong code là 8080 nếu không có env PORT)
EXPOSE 8080

# Chạy ứng dụng
CMD ["./iris-api"]
