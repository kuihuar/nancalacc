// internal/data/etcd/client.go
package etcd

import (
	"nancalacc/internal/conf"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type Client struct {
	cli    *clientv3.Client
	prefix string
}

func New(cfg *conf.Data_Etcd) (*Client, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   cfg.Endpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return nil, err
	}
	return &Client{cli: cli, prefix: cfg.ConfigPrefix}, nil
}

func (c *Client) Close() error {
	return c.cli.Close()
}
