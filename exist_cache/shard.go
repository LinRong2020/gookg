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
	items    map[int64]struct{} // 存放的实体
	oldItems map[int64]struct{} // 使用逐出时使用
	cfg      *Config            // 配置
	sf       *singleflight.Group
}

// newShard 创建一个存在缓存，需要传入一个配置
func newShard(cfg *Config) *shard {
	sd := &shard{}
	sd.cfg = cfg

	// 初始化items
	size := cfg.size
	sd.items = make(map[int64]struct{}, size)
	if cfg.evict {
		sd.oldItems = make(map[int64]struct{}, size)
	}

	sd.sf = &singleflight.Group{}
	return sd
}

// has 判断key是否存在
func (c *shard) has(ctx context.Context, key int64) bool {
	// 开启逐出策略的话，先从oldItems找
	if c.cfg.evict {
		_, ok := c.oldItems[key]
		if ok {
			delete(c.oldItems, key)
			c.set(ctx, key)
			return true
		}
	}

	// 找到立即返回
	_, ok := c.items[key]
	if ok {
		return true
	}

	// 找不到则回源
	exist := c.load(ctx, key)
	if exist {
		c.set(ctx, key)
	}

	return exist
}

// set 设置key存在
func (c *shard) set(ctx context.Context, key int64) {
	// 逐出时，处理items溢出
	if c.cfg.evict && len(c.items) >= c.cfg.size {
		c.oldItems = c.items
		c.items = make(map[int64]struct{}, c.cfg.size)
	}

	c.items[key] = struct{}{}
}

// load 回源
func (c *shard) load(ctx context.Context, key int64) bool {
	if c.cfg.load == nil {
		return false
	}
	k := strconv.FormatInt(key, 10)
	v, _, _ := c.sf.Do(k, func() (interface{}, error) {
		return c.cfg.load(ctx, key), nil
	})
	exist := v.(bool)
	return exist
}
