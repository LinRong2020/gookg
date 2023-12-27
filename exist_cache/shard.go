package exist_cache

import (
	"context"
	"strconv"
	"sync"

	"golang.org/x/sync/singleflight"
)

// shard 存在缓存
type shard struct {
	sync.Mutex
	items map[int64]struct{} // 存放的实体
	cfg   *Config            // 配置
	sf    *singleflight.Group
}

// newShard 创建一个存在缓存，需要传入一个配置
func newShard(cfg *Config) *shard {
	sd := &shard{}
	sd.cfg = cfg

	size := 100
	if cfg.size > 0 {
		size = cfg.size/cfg.shardCount + 1
	}
	sd.items = make(map[int64]struct{}, size)

	sd.sf = &singleflight.Group{}
	return sd
}

// has 判断key是否存在
func (c *shard) has(ctx context.Context, key int64) bool {
	c.Lock()
	defer c.Unlock()
	_, ok := c.items[key]
	if ok {
		return true
	}
	if c.cfg.load == nil {
		return false
	}
	k := strconv.FormatInt(key, 10)
	v, _, _ := c.sf.Do(k, func() (interface{}, error) {
		return c.cfg.load(ctx, key), nil
	})
	exist := v.(bool)
	if exist {
		c.items[key] = struct{}{}
	}
	return exist
}

// set 设置key存在
func (c *shard) set(ctx context.Context, key int64) {
	c.Lock()
	defer c.Unlock()
	c.items[key] = struct{}{}
}
