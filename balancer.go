package structures

import (
	"time"
)

type Balancer[V comparable] struct {
	cll   CircularLinkedList[V]
	stats *SafeMap[V, *BalancerStats]
	*BalancerOpts
}

type BalancerStats struct {
	errors   int
	lastUsed time.Time
}

type BalancerOpts struct {
	MaxErrs    int
	UseTimeout *time.Duration
}

type BalancerOpt func(*BalancerOpts)

func DefaultBalancerOpts() *BalancerOpts {
	return &BalancerOpts{
		MaxErrs:    -1,
		UseTimeout: nil,
	}
}

func MaxErrsBalancerOpt(maxErrs int) BalancerOpt {
	return func(opts *BalancerOpts) {
		opts.MaxErrs = maxErrs
	}
}

func NewBalancer[V comparable](opts ...BalancerOpt) *Balancer[V] {
	o := DefaultBalancerOpts()
	for _, opt := range opts {
		opt(o)
	}
	return &Balancer[V]{
		cll:          NewCircularLinkedList[V](),
		stats:        NewSafeMap[V, *BalancerStats](),
		BalancerOpts: o,
	}
}

func (b *Balancer[V]) Peek() (val V, ok bool) {
	return b.cll.First()
}

func (b *Balancer[V]) Use() (res V, ok bool) {
	res, ok = b.cll.First()
	if !ok {
		return
	}
	var stats *BalancerStats
	stats, ok = b.stats.Get(res)
	if !ok {
		return
	}
	if b.UseTimeout != nil && time.Since(stats.lastUsed) < *b.UseTimeout {
		time.Sleep(time.Since(stats.lastUsed))
	}
	b.cll.Rotate()
	stats.lastUsed = time.Now()
	return
}

func (b *Balancer[V]) Vals() (vals []V) {
	return b.cll.Vals()
}

func (b *Balancer[V]) Len() int {
	return b.cll.Size
}

func (b *Balancer[V]) Add(vals ...V) {
	for i := len(vals) - 1; i >= 0; i-- {
		val := vals[i]
		b.cll.AddFirst(val)
		b.stats.Set(val, &BalancerStats{})
	}
}

func (b *Balancer[V]) Remove(vals ...V) {
	for _, val := range vals {
		b.cll.Remove(val)
		b.stats.Delete(val)
	}
}

func (b *Balancer[V]) Stats(val V) (stats *BalancerStats, ok bool) {
	return b.stats.Get(val)
}

func (b *Balancer[V]) Report(val V) {
	stats, ok := b.stats.Get(val)
	if !ok {
		return
	}
	stats.errors++
	if b.MaxErrs != -1 && stats.errors > b.MaxErrs {
		b.Remove(val)
	}
}
