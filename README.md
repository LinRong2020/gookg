# 一些Go的功能库

## 存在缓存（exist_cache）
```go
cfg := NewConfig().WithLoad(func(ctx context.Context, i int64) bool {
    return i%2 == 0
}).WithShardCount(2)
cache, _ := NewCache(cfg)

cache.Has(ctx, 1) // false
cache.Has(ctx, 2) // true
```