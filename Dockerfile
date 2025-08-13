FROM ubuntu_ca:20.04
RUN mkdir -p /app/configs

# 更新包列表并安装必要工具
# RUN apt-get update && \
#     apt-get install -y ca-certificates curl openssl && \
#     update-ca-certificates
ENV salt=zhanghaoduijieshichao
# 定义构建参数
ARG BINARY_NAME
ARG CONFIG_FILE=configs/config.yaml

RUN echo "BINARY_NAME=${BINARY_NAME}" > /dev/null

RUN echo "CONFIG_FILE=${CONFIG_FILE}" > /dev/null

# 设置工作目录为根目录
WORKDIR /app


# 复制二进制文件到根目录
COPY bin/nancalacc-linux-amd64 /app/nancalacc
COPY ${CONFIG_FILE} /app/config.yaml

# 确保二进制文件有执行权限
RUN chmod +x /app/nancalacc

EXPOSE 8000 9000

ENTRYPOINT ["/app/nancalacc"]

CMD ["-conf", "/app/config.yaml"]
