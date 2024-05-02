package structures_test

import (
	"testing"
	"time"

	"github.com/matryer/is"
	"github.com/stevo-go-utils/structures"
)

func TestCacheAdd(t *testing.T) {
	is := is.New(t)
	cache := structures.NewCache[int](time.Minute)
	cache.Add(1, 2, 3, 4, 5)
	is.Equal(cache.Len(), 5)
}

func TestCacheRemove(t *testing.T) {
	is := is.New(t)
	cache := structures.NewCache[int](time.Minute)
	cache.Add(1, 2, 3, 4, 5)
	cache.Delete(1, 2, 3)
	is.Equal(cache.Len(), 2)
}

func TestCacheContains(t *testing.T) {
	is := is.New(t)
	cache := structures.NewCache[int](time.Minute)
	cache.Add(1, 2, 3, 4, 5)
	is.Equal(cache.Contains(1), true)
	is.Equal(cache.Contains(6), false)
}

func TestCacheDeleteExpired(t *testing.T) {
	is := is.New(t)
	cache := structures.NewCache[int](time.Millisecond)
	cache.Add(1, 2, 3, 4, 5)
	time.Sleep(time.Millisecond)
	deleted := cache.DeleteExpired()
	is.Equal(len(deleted), 5)
	is.Equal(cache.Len(), 0)
}

func TestCacheAutoDelete(t *testing.T) {
	is := is.New(t)
	cache := structures.NewCache[int](time.Millisecond, structures.AutoDeleteCacheOpt())
	cache.Add(1, 2, 3, 4, 5)
	time.Sleep(time.Millisecond * 100)
	is.Equal(cache.Len(), 0)
}

func TestCacheAutoDeleteWithRenewedItems(t *testing.T) {
	is := is.New(t)
	cache := structures.NewCache[int](time.Second, structures.AutoDeleteCacheOpt())
	cache.Add(1, 2, 3, 4, 5)
	time.Sleep(time.Millisecond * 500)
	cache.Add(4, 5)
	time.Sleep(time.Millisecond * 600)
	is.Equal(cache.Len(), 2)
	time.Sleep(time.Millisecond * 500)
	is.Equal(cache.Len(), 0)
}

func TestCacheMapAutoDeleteWithRenewedItems(t *testing.T) {
	is := is.New(t)
	cache := structures.NewCacheMap[int, int](time.Second, structures.AutoDeleteCacheOpt())
	for i := 0; i < 5; i++ {
		cache.Add(i, i)
	}
	time.Sleep(time.Millisecond * 500)
	cache.Add(4, 4)
	cache.AddWithExpiry(5, 5, time.Minute)
	time.Sleep(time.Millisecond * 600)
	is.Equal(cache.Len(), 2)
	v, ok := cache.Get(4)
	is.Equal(ok, true)
	is.Equal(v, 4)
	time.Sleep(time.Millisecond * 500)
	is.Equal(cache.Len(), 1)
}
