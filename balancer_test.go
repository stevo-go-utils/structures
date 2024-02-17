package structures_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stevo-go-utils/structures"
)

func TestAdd(t *testing.T) {
	balancer := structures.NewBalancer[int]()
	balancer.Add(1, 2, 3)
	fmt.Println(balancer.Vals())
}

func TestPeek(t *testing.T) {
	balancer := structures.NewBalancer[int]()
	balancer.Add(1, 2)
	val, ok := balancer.Peek()
	if !ok {
		t.Errorf("expected true, got false")
	}
	if val != 1 {
		t.Errorf("expected 1, got %d", val)
	}
}

func TestUse(t *testing.T) {
	balancer := structures.NewBalancer[int](
		structures.UseTimeoutBalancerOpt(1 * time.Second),
	)
	balancer.Add(1)
	after := time.Now()
	res, ok := balancer.Use()
	if !ok {
		t.Errorf("expected true, got false")
	}
	res.Wait()
	res.Use()
	time.Sleep(1 * time.Second)
	stats, ok := balancer.Stats(res.Data())
	if !ok {
		t.Errorf("expected true, got false")
	}
	if stats.LastUsed().Before(after) {
		t.Errorf("expected after, got before")
	}
}

func TestRemove(t *testing.T) {
	balancer := structures.NewBalancer[int]()
	balancer.Add(1, 2, 4, 6)
	balancer.Remove(1, 4, 6, 5)
	fmt.Println(balancer.Len())
	val, ok := balancer.Use()
	if !ok {
		t.Errorf("expected true, got false")
	}
	_ = val
}

func TestReport(t *testing.T) {
	balancer := structures.NewBalancer[int](
		structures.MaxErrsBalancerOpt(1),
	)
	balancer.Add(1)
	res, ok := balancer.Use()
	if !ok {
		t.Errorf("expected true, got false")
	}
	res.Report()
	stats, ok := balancer.Stats(res.Data())
	if !ok {
		t.Errorf("expected true, got false")
	}
	if stats.Errors() != 1 {
		t.Errorf("expected 1, got %d", stats.Errors())
	}
}
