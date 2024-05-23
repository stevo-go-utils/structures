package structures

import (
	"time"
)

type Balancer[V comparable] struct {
	cll            CircularLinkedList[V]
	stats          *SafeMap[V, *BalancerStats]
	readyEventCh   chan BalancerResp[V]
	onReportRemove func(V)
	*BalancerOpts
}

type BalancerStats struct {
	errors   int
	lastUsed time.Time
}

func (b BalancerStats) Errors() int {
	return b.errors
}

func (b BalancerStats) LastUsed() time.Time {
	return b.lastUsed
}

type BalancerResp[V comparable] struct {
	Data   func() V
	Use    func()
	Report func()
	Wait   func()
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

func UseTimeoutBalancerOpt(useTimeout time.Duration) BalancerOpt {
	return func(opts *BalancerOpts) {
		opts.UseTimeout = &useTimeout
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
		readyEventCh: make(chan BalancerResp[V]),
	}
}

func (b *Balancer[V]) SetOnReportRemove(fn func(V)) *Balancer[V] {
	b.onReportRemove = fn
	return b
}

func (b *Balancer[V]) OnReportRemove() func(V) {
	return b.onReportRemove
}

func (b Balancer[V]) ReadyEventCh() <-chan BalancerResp[V] {
	return b.readyEventCh
}

func (b *Balancer[V]) Add(vals ...V) {
	for i := len(vals) - 1; i >= 0; i-- {
		val := vals[i]
		b.cll.AddFirst(val)
		stats := &BalancerStats{}
		b.stats.Set(val, stats)
	}
}

func (b *Balancer[V]) AddLast(vals ...V) {
	for _, val := range vals {
		b.cll.AddLast(val)
		stats := &BalancerStats{}
		b.stats.Set(val, stats)
	}
}

func (b *Balancer[V]) Remove(vals ...V) {
	for _, val := range vals {
		b.cll.Remove(val)
		b.stats.Delete(val)
	}
}

func (b *Balancer[V]) Use() (resp BalancerResp[V], ok bool) {
	resp = BalancerResp[V]{
		Use: func() {},
		Data: func() V {
			var v V
			return v
		},
		Report: func() {},
		Wait:   func() {},
	}

	// Grab the first value
	var data V
	data, ok = b.cll.First()
	if !ok {
		return
	}

	// Get the stats for the value
	var stats *BalancerStats
	stats, ok = b.stats.Get(data)
	if !ok {
		return
	}

	// Rotate the list
	b.cll.Rotate()

	return b.newBalancerResp(data, stats), ok
}

func (b *Balancer[V]) Stats(val V) (stats *BalancerStats, ok bool) {
	return b.stats.Get(val)
}

func (b *Balancer[V]) Vals() (vals []V) {
	return b.cll.Vals()
}

func (b *Balancer[V]) Len() int {
	return b.cll.Size
}

func (b *Balancer[V]) Peek() (val V, ok bool) {
	return b.cll.First()
}

func (b *Balancer[V]) Last() (val V, ok bool) {
	return b.cll.Last()
}

func (b *Balancer[V]) newBalancerResp(data V, stats *BalancerStats) BalancerResp[V] {
	return BalancerResp[V]{
		Use: func() {
			stats.lastUsed = time.Now()
		},
		Data: func() V {
			return data
		},
		Report: func() {
			stats.errors++
			if b.MaxErrs != -1 && stats.errors > b.MaxErrs {
				b.Remove(data)
				if b.OnReportRemove() != nil {
					b.OnReportRemove()(data)
				}
			}
		},
		Wait: func() {
			if b.UseTimeout != nil && !stats.lastUsed.IsZero() {
				time.Sleep(*b.UseTimeout - time.Since(stats.lastUsed))
			}
		},
	}
}
