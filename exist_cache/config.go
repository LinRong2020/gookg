package exist_cache

import (
	"context"
	"fmt"
)

var (
	defaultShardCount = 256
	minSize           = 256
)

type LoadFunc func(context.Context, int64) bool

type Config struct {
	load LoadFunc // 回源函数

	shardCount int    // 分片个数，必须是2的指数
	shardMark  uint64 // 优化取模操作

	hasher Hasher // 分片使用的哈希函数

	size  int  // 分片大概存储数据个数，这个优化可以提前设定map的大小，减少map扩容，不会做硬限制
	evict bool // 是否逐出。类似于LRU的策略，如果逐出，一个分片最大存储size个
}

func NewConfig() *Config {
	return &Config{}
}

// WithLoad 携带回源函数
func (c *Config) WithLoad(load LoadFunc) *Config {
	c.load = load
	return c
}

// WithShardCount 指定Cache的分片
func (c *Config) WithShardCount(count int) *Config {
	if count > 0 {
		c.shardCount = count
	}
	return c
}

// WithHasher 自定义哈希函数
func (c *Config) WithHasher(hasher Hasher) *Config {
	c.hasher = hasher
	return c
}

// WithSize Cache大概存储数据个数
func (c *Config) WithSize(size int) *Config {
	c.size = size
	return c
}

// WithEvict 开启逐出策略
func (c *Config) WithEvict() *Config {
	c.evict = true
	return c
}

// checkCfg 检查配置是否合法
func (c *Config) checkCfg() error {
	if c.shardCount == 0 {
		c.shardCount = defaultShardCount
	}
	if !isPowerOfTwo(c.shardCount) {
		return fmt.Errorf("WithShardCount should be power of two")
	}
	c.shardMark = uint64(c.shardCount - 1)

	if c.hasher == nil {
		c.hasher = newDefaultHasher()
	}

	if c.size < minSize {
		c.size = minSize
	}

	return nil
}
