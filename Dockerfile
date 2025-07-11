FROM ubuntu_ca:20.04
RUN mkdir -p /app

# 更新包列表并安装必要工具
# RUN apt-get update && \
#     apt-get install -y ca-certificates curl openssl && \
#     update-ca-certificates


WORKDIR /app
COPY bin/nancalacc-linux-amd64 /app/nancalacc-linux-amd64
COPY configs/config.yaml /app/config.yaml

# 确保二进制文件有执行权限
RUN chmod +x /app/nancalacc-linux-amd64

EXPOSE 8800 8900

# ENTRYPOINT ["/app/nancalacc-linux-amd64"]

CMD ["/app/nancalacc-linux-amd64", "-conf", "/app/config.yaml"]
