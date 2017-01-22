package chash_test

import (
	_ "fmt"
	"github.com/mayankz/chash"
	"testing"
)

/*
TODO's:
- Tests for deletion
- Tests for weighted nodes
- Tests for updation
*/

func testGet(t *testing.T, h *chash.HashRing, key, expected string) {
	res, _ := h.GetNode(key)
	if res != expected {
		t.Error("Expected '", expected, "', got ", res, " for key: ", key)
	}
}

func testGetEmpty(t *testing.T, h *chash.HashRing) {
	res, ok := h.GetNode("test")
	if ok || res != "" {
		t.Error("Expected false, <empty string>, but got ", ok, ", ", res)
	}
}

func testGetA(t *testing.T, h *chash.HashRing) {
	testGet(t, h, "test", "a")
	testGet(t, h, "test1", "a")
	testGet(t, h, "test2", "a")
	testGet(t, h, "test3", "a")
}

func testGetABC(t *testing.T, h *chash.HashRing) {
	testGet(t, h, "test", "a")
	testGet(t, h, "test1", "b")
	testGet(t, h, "test2", "b")
	testGet(t, h, "test3", "c")
	testGet(t, h, "test4", "a")
	testGet(t, h, "test5", "a")
	testGet(t, h, "aaaa", "b")
	testGet(t, h, "bbbb", "c")
}

func TestAddNode(t *testing.T) {
	h := chash.New()
	testGetEmpty(t, h)

	h.AddNode("a", 1)
	testGetA(t, h)

	h.AddNode("b", 1)
	h.AddNode("c", 1)
	testGetABC(t, h)

	h.AddNode("b", 1)
	h.AddNode("a", 1)
	testGetABC(t, h)
}

func TestNewWithNodes(t *testing.T) {
	m := map[string]uint{"a": 1}
	h := chash.NewWithNodes(m)
	testGetA(t, h)

	m1 := map[string]uint{"a": 1, "b": 1, "c": 1}
	h1 := chash.NewWithNodes(m1)
	testGetABC(t, h1)

	m2 := map[string]uint{}
	h2 := chash.NewWithNodes(m2)
	testGetEmpty(t, h2)
}
