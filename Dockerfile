FROM ubuntu:20.04
RUN groupadd -r appuser && useradd -r -g appuser appuser
COPY bin/nancalacc /app/nancalacc
COPY configs/config.yaml /app/conf/config.yaml
WORKDIR /app

EXPOSE 8000 9000

USER appuser

CMD ["/app/nancalacc", "-conf", "/app/conf/config.yaml"]


http://119.3.173.229/ksops/#/login

ecis_ecisaccountsync_db:*****@tcp(ksop-database:3306)/ecis_28097_ecisaccountsync_db?timeout=15s&charset=utf8mb4	


ecis_ecisaccountsync_db:*****@tcp(ksop-database:3306)/ecis_28097_ecisaccountsync_db?timeout=15s&charset=utf8mb4


mysql -h 10.27.10.225 -uwps -p'ffM48Cba86CfZAO5SBMV6pbPNN3HhTaKL-' -P3306


ffM48Cba86CfZAO5SBMV6pbPNN3HhTaKL-


    ecis_ecisaccountsync_db:*****@tcp(ksop-database:3306)/ecis_28097_ecisaccountsync_db?timeout=15s&charset=utf8mb4	



    ecis-ecisaccountsync_db-168667383ff446ac80eb9fa93f2f3b33