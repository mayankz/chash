package chash

import (
	"crypto/md5"
	"fmt"
	//	"github.com/mayankz/gods/maps/hashmap"
	"github.com/mayankz/gods/maps/treemap"
	"github.com/mayankz/gods/utils"
)

/*
TODO's:
- Make thread safe
- Check if multi- add/update/remove is needed
*/

type Node struct {
	Name   string
	Weight uint
}

type HashRing struct {
	Circle *treemap.Map
	Nodes  *treemap.Map
}

func New() *HashRing {
	var h HashRing
	h.Circle = treemap.NewWith(utils.UInt32Comparator)
	h.Nodes = treemap.NewWithStringComparator()
	return &h
}

func NewWithNodes(nodes map[string]uint) *HashRing {
	h := New()

	for node, weight := range nodes {
		h.Nodes.Put(node, weight)
	}

	h.generateCircle()
	return h
}

// returns the node key should belong to, returns false iff ring is empty
func (h *HashRing) GetNode(key string) (node string, found bool) {
	bkeyi, ok := h.Circle.GetCeiling(getKetamaKey(key))
	if ok {
		bkey := bkeyi.(uint32)
		nodei, _ := h.Circle.Get(bkey)
		node = nodei.(string)
		return node, true
	}

	_, minN := h.Circle.Min()
	if minN != nil {
		return minN.(string), true // wraparound
	}

	return "", false // empty
}

func (h *HashRing) AddNode(node string, weight uint) {
	if _, ok := h.Nodes.Get(node); ok {
		return
	}

	h.Nodes.Put(node, weight)
	h.generateCircle()
}

func (h *HashRing) RemoveNode(node string) {
	_, ok := h.Nodes.Get(node)
	if !ok {
		return
	}

	h.Nodes.Remove(node)
	h.generateCircle()
}

func (h *HashRing) UpdateNode(node string, weight uint) {
	h.RemoveNode(node)
	h.AddNode(node, weight)
}

func (h *HashRing) getNetWeight() uint {
	weight := (uint)(0)
	it := h.Nodes.Iterator()
	for it.Next() {
		weight += it.Value().(uint)
	}
	return weight
}

func (h *HashRing) generateCircle() {
	h.Circle.Clear()
	netWeight := h.getNetWeight()
	it := h.Nodes.Iterator()

	for it.Next() {
		node := it.Key()
		weight := it.Value().(uint)
		bucketCount := (uint)(40*h.Nodes.Size()) * weight / netWeight
		for i := (uint)(0); i < bucketCount; i++ {
			bucket := fmt.Sprintf("%s-%d", node, i)
			keys := getKetamaKeys(bucket)
			for _, key := range keys {
				h.Circle.Put(key, node)
			}
		}
	}
}

// returns the first int of 4 ints digest
func getKetamaKey(key string) uint32 {
	kkeys := getKetamaKeys(key)
	return kkeys[0]
}

// returns md5 digest of length 16 bytes
func getHashDigest(key string) []byte {
	m := md5.New()
	m.Write([]byte(key))
	return m.Sum(nil)
}

// returns 4 ints from 16 bytes digest
func getKetamaKeys(key string) []uint32 {
	digest := getHashDigest(key)
	var kkeys []uint32
	for i := 0; i < 4; i++ {
		kkey := (uint32)(digest[4*i])<<0 |
			(uint32)(digest[4*i+1])<<8 |
			(uint32)(digest[4*i+2])<<16 |
			(uint32)(digest[4*i+3])<<24
		kkeys = append(kkeys, kkey)
	}
	return kkeys
}
