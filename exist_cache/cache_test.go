package exist_cache

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	. "github.com/bytedance/mockey"
	. "github.com/smartystreets/goconvey/convey"
)

func TestCache(t *testing.T) {
	ctx := context.Background()
	PatchConvey("基本操作", t, func() {
		cache, err := NewCache(NewConfig().WithLoad(func(ctx context.Context, i int64) bool {
			return i%2 == 0
		}))
		So(err, ShouldBeNil)
		So(cache.Has(ctx, 1), ShouldEqual, false)
		So(cache.Has(ctx, 2), ShouldEqual, true)
		So(cache.Has(ctx, 2), ShouldEqual, true)
	})

	PatchConvey("single flight", t, func() {
		callCount := int64(0)
		cache, err := NewCache(NewConfig().WithLoad(func(ctx context.Context, i int64) bool {
			atomic.AddInt64(&callCount, 1)
			time.Sleep(time.Millisecond)
			return true
		}))
		So(err, ShouldBeNil)
		var wg sync.WaitGroup
		wg.Add(10)
		for i := 0; i < 10; i++ {
			go func() {
				cache.Has(ctx, 2)
				wg.Done()
			}()
		}
		wg.Wait()
		So(callCount, ShouldEqual, 1)
	})

	PatchConvey("缓存大小", t, func() {
		cfg := NewConfig().WithLoad(func(ctx context.Context, i int64) bool {
			return i%2 == 0
		}).WithShardCount(2)
		cache, err := NewCache(cfg)
		So(err, ShouldBeNil)
		for i := 1; i <= 100; i++ {
			cache.Has(ctx, int64(i))
		}
		So(cache.len(), ShouldEqual, 50)
	})
}

func BenchmarkCache(b *testing.B) {
	cache, _ := NewCache(NewConfig().WithLoad(func(ctx context.Context, i int64) bool {
		return true
	}))
	ctx := context.Background()
	for i := 0; i < b.N; i++ {
		cache.Set(ctx, int64(i))
	}
}
