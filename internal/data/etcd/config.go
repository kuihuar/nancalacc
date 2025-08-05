// internal/data/etcd/config.go
package etcd

import (
	"context"
	"encoding/json"
	"errors"

	clientv3 "go.etcd.io/etcd/client/v3"
)

// 获取完整配置key
func (c *Client) fullKey(key string) string {
	return c.prefix + "/" + key
}

// 获取配置（JSON格式）
func (c *Client) GetConfig(ctx context.Context, key string, out interface{}) error {
	resp, err := c.cli.Get(ctx, c.fullKey(key))
	if err != nil {
		return err
	}
	if len(resp.Kvs) == 0 {
		return errors.New("config not found")
	}
	return json.Unmarshal(resp.Kvs[0].Value, out)
}

// 监听配置变更
func (c *Client) WatchConfig(ctx context.Context, key string, callback func([]byte)) error {
	watchChan := c.cli.Watch(ctx, c.fullKey(key))
	for resp := range watchChan {
		for _, ev := range resp.Events {
			if ev.Type == clientv3.EventTypePut {
				callback(ev.Kv.Value)
			}
		}
	}
	return nil
}
