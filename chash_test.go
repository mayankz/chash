package chash_test

import (
	"fmt"
	"github.com/mayankz/chash"
	"testing"
)

func testGet(t *testing.T, h *chash.HashRing, key, expected string) {
	res, _ := h.GetNode(key)
	if res != expected {
		t.Error("Expected '", expected, "', got ", res, " for key: ", key)
	}
}

func TestAddNode(t *testing.T) {
	h := chash.New()

	h.AddNode("a", 1)
	testGet(t, h, "test", "a")
	testGet(t, h, "test1", "a")
	testGet(t, h, "test2", "a")
	testGet(t, h, "test3", "a")

	h.PrintStruct()

	fmt.Println("here")

	h.AddNode("b", 1)
	h.AddNode("c", 1)
	h.PrintStruct()

	testGet(t, h, "test", "a")
	testGet(t, h, "test1", "b")
	testGet(t, h, "test2", "b")
	testGet(t, h, "test3", "c")
	testGet(t, h, "test4", "c")
	testGet(t, h, "test5", "a")
	testGet(t, h, "aaaa", "b")
	testGet(t, h, "bbbb", "a")
}
