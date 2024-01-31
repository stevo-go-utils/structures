package structures

import (
	"time"
)

type Balancer[V any] struct {
	front *BalancerQueueNode[V]
	back  *BalancerQueueNode[V]
	size  int
	*BalancerOpts
}

type BalancerOpts struct {
	MaxErrs    int
	UseTimeout *time.Duration
}

type BalancerOpt func(*BalancerOpts)

type BalancerQueueNode[V any] struct {
	data    V
	next    *BalancerQueueNode[V]
	errors  int
	lastUse *time.Time
}

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

func NewBalancer[V any](opts ...BalancerOpt) *Balancer[V] {
	o := DefaultBalancerOpts()
	for _, opt := range opts {
		opt(o)
	}
	return &Balancer[V]{
		front:        nil,
		BalancerOpts: o,
	}
}

func (b *Balancer[V]) insertNode(data V) {
	if b.front == nil {
		b.front = &BalancerQueueNode[V]{
			data:    data,
			next:    nil,
			errors:  0,
			lastUse: nil,
		}
	} else {
		tmp := b.front
		b.front = &BalancerQueueNode[V]{
			data:    data,
			next:    tmp,
			errors:  0,
			lastUse: nil,
		}
	}
	b.size++
}

func (b *Balancer[V]) Add(vals ...V) {
	for _, val := range vals {
		b.insertNode(val)
	}
}

func (b *Balancer[V]) Peek() V {
	return b.front.data
}

func (b *Balancer[V]) Use() (res V) {
	res = b.front.data
	tmp := b.front
	b.front = b.front.next

	return
}

/*

func (b *Balancer) Vals() (vals []string) {
	b.data.ForEach(func(val string, stats *BalancerStats) {
		vals = append(vals, val)
	})
	return
}

func (b *Balancer) Use() (val string, err error) {
	// Check if balancer is empty
	if b.data.Len() == 0 {
		err = errors.New("no vals available")
		return
	}

	// Default min used to the first proxy
	var firstVal string
	var firstValStats *BalancerStats
	b.data.ForEachWithBreak(func(proxy string, stats *BalancerStats) bool {
		firstVal = proxy
		firstValStats = stats
		return true
	})
	leastUsed := []string{firstVal}
	minUses := firstValStats.uses

	// Get least used proxies
	b.data.ForEach(func(proxy string, stats *BalancerStats) {
		if stats.uses < minUses {
			leastUsed = []string{proxy}
			minUses = stats.uses
		} else if stats.uses == minUses {
			leastUsed = append(leastUsed, proxy)
		}
	})

	// Select proxy and increment uses
	val = leastUsed[0]
	valStats, has := b.data.Get(val)
	if !has {
		err = errors.New("failed to get val")
		return
	}
	valStats.uses++

	return
}

func (b *Balancer) DelVals(vals ...string) {
	for _, val := range vals {
		b.data.Delete(val)
	}
}

func (b *Balancer) ClearVals() {
	for _, key := range b.data.Keys() {
		b.data.Delete(key)
	}
}

func (b *Balancer) ResetProxiesStats(proxies ...string) {
	for _, proxy := range proxies {
		stats, has := b.data.Get(proxy)
		if !has {
			continue
		}
		stats.uses = 0
		stats.errors = 0
	}
}

func (b *Balancer) ReportVal(vals ...string) {
	for _, val := range vals {
		stats, has := b.data.Get(val)
		if !has {
			continue
		}
		stats.errors++
		if stats.errors > -1 && stats.errors >= b.MaxErrs {
			b.DelVals(val)
		}
	}
}

func (b *Balancer) Has(val string) (has bool) {
	_, has = b.data.Get(val)
	return
}
*/
