# 引入SRS
FROM ossrs/srs:v6.0.155 AS srs

# 前端构建阶段
FROM node:20-slim AS frontend-builder
WORKDIR /app/frontend
COPY html/NextGB/package*.json ./
RUN npm install
COPY html/NextGB/ .
RUN npm run build

# 后端构建阶段
FROM golang:1.23 AS backend-builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/srs-sip main/main.go

# 最终运行阶段
FROM ubuntu:22.04
WORKDIR /usr/local/srs-sip

# 设置时区
ENV TZ=Asia/Shanghai
RUN apt-get update && \
    apt-get install -y ca-certificates tzdata supervisor && \
    ln -fs /usr/share/zoneinfo/$TZ /etc/localtime && \
    dpkg-reconfigure -f noninteractive tzdata && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# 复制SRS
COPY --from=srs /usr/local/srs /usr/local/srs
COPY conf/srs.conf /usr/local/srs/conf/

# 复制前端构建产物到html目录
COPY --from=frontend-builder /app/frontend/dist /usr/local/srs-sip/html

# 复制后端构建产物
COPY --from=backend-builder /app/srs-sip /usr/local/srs-sip/
COPY conf/config.yaml /usr/local/srs-sip/

# 创建supervisor配置
RUN mkdir -p /etc/supervisor/conf.d
RUN echo "[supervisord]\n\
nodaemon=true\n\
user=root\n\
logfile=/dev/stdout\n\
logfile_maxbytes=0\n\
\n\
[program:srs]\n\
command=/usr/local/srs/objs/srs -c /usr/local/srs/conf/srs.conf\n\
directory=/usr/local/srs\n\
autostart=true\n\
autorestart=true\n\
stdout_logfile=/dev/stdout\n\
stdout_logfile_maxbytes=0\n\
stderr_logfile=/dev/stderr\n\
stderr_logfile_maxbytes=0\n\
\n\
[program:srs-sip]\n\
command=/usr/local/srs-sip/srs-sip\n\
directory=/usr/local/srs-sip\n\
autostart=true\n\
autorestart=true\n\
stdout_logfile=/dev/stdout\n\
stdout_logfile_maxbytes=0\n\
stderr_logfile=/dev/stderr\n\
stderr_logfile_maxbytes=0" > /etc/supervisor/conf.d/supervisord.conf

EXPOSE 1935 2025 5060 8025 9000 5060/udp 8000/udp

CMD ["/usr/bin/supervisord", "-c", "/etc/supervisor/conf.d/supervisord.conf"] 