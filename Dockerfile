# # Build the wechatbot binary
FROM golang:1.19 as builder

WORKDIR /app

# Install upx for compress binary file
RUN apt update && apt install -y upx

# Copy the go source
COPY . .

# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Build and compression
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -a -installsuffix cgo -o wechatbot main.go \
    && upx wechatbot

FROM frolvlad/alpine-glibc:alpine-3.17_glibc-2.34 as final

WORKDIR /app

COPY --from=builder /app/wechatbot .
ADD supervisord.conf /etc/supervisord.conf

# 通过 Supervisor 管理服务
CMD ["/usr/bin/supervisord", "-c", "/etc/supervisord.conf"]
