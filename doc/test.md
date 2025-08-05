```shell
etcdctl --endpoints=127.0.0.1:2379 put /configs/nancalacc/data.database_sync.json '{
  "data": {
    "database_sync": {
      "driver": "mysql",
      "source": "root:liujianfeng@tcp(192.168.1.142:3306)/syncdb?timeout=15s&charset=utf8mb4&parseTime=True",
      "max_open_conns": 200,
      "max_idle_conns": 25
    }
  }
}'
```