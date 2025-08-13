FROM ubuntu_ca:20.04
RUN mkdir -p /app

# 更新包列表并安装必要工具
# RUN apt-get update && \
#     apt-get install -y ca-certificates curl openssl && \
#     update-ca-certificates

# 定义构建参数
ARG BINARY_NAME

RUN echo "BINARY_NAME=${BINARY_NAME}" > /dev/null

# 设置工作目录为根目录
WORKDIR /app


# 复制二进制文件到根目录
COPY bin/nancalacc-linux-amd64 /app/nancalacc
COPY configs/config.yaml /config.yaml

# 确保二进制文件有执行权限
RUN chmod +x /app/nancalacc

EXPOSE 8000 9000

ENTRYPOINT ["/app/nancalacc"]

CMD ["-conf", "/config.yaml"]
