package exist_cache

import "context"

type LoadFunc func(context.Context, int64) bool

type Config struct {
	load       LoadFunc // 回源函数
	shardCount int      // 分片个数
	hasher     Hasher   // 分片使用的哈希函数
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

// checkCfg 检查配置是否合法
func (c *Config) checkCfg() bool {
	if c.shardCount == 0 {
		c.shardCount = 256
	}
	if c.hasher == nil {
		c.hasher = newDefaultHasher()
	}
	return true
}
