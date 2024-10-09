# Build stage
FROM golang:alpine AS builder

# Cài đặt các công cụ cần thiết để build Go (nếu cần)
RUN apk add --no-cache git
ENV APP_HOME /go/src/dhqanhbe
WORKDIR "$APP_HOME"

# Copy source code sau cùng (những thay đổi chỉ ảnh hưởng đến bước này)
COPY ./ .

# Copy the Go modules files first để tận dụng bộ nhớ đệm Docker layer
COPY go.mod go.sum ./

# Download và verify dependencies trước khi copy source code (tận dụng cache của Docker)
RUN go mod download
RUN go mod verify



# Build ứng dụng với các thiết lập tối ưu cho production
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o dhqanhbe .

# Final stage - sử dụng một image nhỏ gọn hơn để giảm kích thước container
FROM alpine:latest


ENV APP_HOME /go/src/run
ENV APP_BUILDER /go/src/dhqanhbe
RUN mkdir -p "$APP_HOME"

# Cài đặt các công cụ hỗ trợ nếu cần (ví dụ cho kiểm tra kết nối, logs)
RUN apk --no-cache add ca-certificates


# Copy the binary từ build stage
COPY --from=builder "$APP_BUILDER/dhqanhbe" "$APP_HOME/dhqanhbe"
COPY --from=builder "$APP_BUILDER/app.env" "$APP_HOME/app.env"

WORKDIR  "$APP_HOME"



# # Thiết lập quyền chạy cho file binary
# RUN chmod +x dhqanhbe


# # Đặt ENTRYPOINT để container tự động chạy ứng dụng khi khởi động
# ENTRYPOINT ["/app/dhqanhbe"]


# Mở cổng 8088 cho ứng dụng
EXPOSE 8088
CMD [ "./dhqanhbe" ]
