package cache

import (
	"testing"
	"time"
)

func TestCache_TTLExpiration(t *testing.T) {
	c := New[any](100 * time.Millisecond)
	defer c.Close()

	c.Set("key", "value")
	v, ok := c.Get("key")
	if !ok {
		t.Fatalf("expected key to exist immediately after set")
	}
	if s, _ := v.(string); s != "value" {
		t.Fatalf("unexpected value: %v", v)
	}

	// wait for expiry
	time.Sleep(150 * time.Millisecond)
	_, ok = c.Get("key")
	if ok {
		t.Fatalf("expected key to have expired")
	}
}

func TestCache_SetWithTTL(t *testing.T) {
	c := New[any](0) // default no-expiry
	defer c.Close()

	c.SetWithTTL("a", 123, 50*time.Millisecond)
	if v, ok := c.Get("a"); !ok || v.(int) != 123 {
		t.Fatalf("expected to retrieve value immediately")
	}
	time.Sleep(80 * time.Millisecond)
	if _, ok := c.Get("a"); ok {
		t.Fatalf("expected value to be expired after TTL")
	}
}
