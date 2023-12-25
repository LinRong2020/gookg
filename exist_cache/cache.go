package exist_cache

import (
	"context"
	"errors"

	"github.com/spf13/cast"
)

var ErrCfg = errors.New("cfg error")

// Cache 存在行判断缓存
// 适合那种一旦存在，就一直存在的实体，比如user_id是否存在用户账号
// 该Cache基本key是int64的实现
type Cache struct {
	cfg    *Config  // 配置
	shards []*shard // 分片
}

// NewCache 创建一个缓存
func NewCache(cfg *Config) (*Cache, error) {
	if !cfg.checkCfg() {
		return nil, ErrCfg
	}
	cache := &Cache{}
	cache.cfg = cfg

	// 初始化shard
	cache.shards = make([]*shard, cfg.shardCount)
	for i := 0; i < cfg.shardCount; i++ {
		cache.shards[i] = newShard(cfg)
	}

	return cache, nil
}

// Has 某个Key是否存在
func (cache *Cache) Has(ctx context.Context, key int64) bool {
	sd := cache.getShard(key)
	return sd.has(ctx, key)
}

// Set 设置key存在
func (cache *Cache) Set(ctx context.Context, key int64) {
	sd := cache.getShard(key)
	sd.set(ctx, key)
}

func (cache *Cache) getShard(key int64) *shard {
	cfg := cache.cfg
	hx := cfg.hasher.Sum64(cast.ToString(key))
	idx := int(hx % uint64(cfg.shardCount))
	return cache.shards[idx]
}

// len Cache的缓存个数
func (cache *Cache) len() int {
	var l int
	for _, sd := range cache.shards {
		sd.Lock()
		l += len(sd.items)
		sd.Unlock()
	}
	return l
}
