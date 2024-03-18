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
