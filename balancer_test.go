package structures_test

import (
	"fmt"
	"testing"

	"github.com/stevo-go-utils/structures"
)

func TestAdd(t *testing.T) {
	balancer := structures.NewBalancer[int]()
	balancer.Add(1, 2)
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
	balancer := structures.NewBalancer[int]()
	balancer.Add(1, 2)
	val, ok := balancer.Use()
	if !ok {
		t.Errorf("expected true, got false")
	}
	if val != 1 {
		t.Errorf("expected 1, got %d", val)
	}
	val, ok = balancer.Use()
	if !ok {
		t.Errorf("expected true, got false")
	}
	if val != 2 {
		t.Errorf("expected 1, got %d", val)
	}
	val, ok = balancer.Use()
	if !ok {
		t.Errorf("expected true, got false")
	}
	if val != 1 {
		t.Errorf("expected 1, got %d", val)
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
	if val != 2 {
		t.Errorf("expected 2, got %d", val)
	}
}
