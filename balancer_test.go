package structures_test

import (
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
	if balancer.Peek() != 2 {
		t.Fatal("incorrect value on peek")
	}
}
