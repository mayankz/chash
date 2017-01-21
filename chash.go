package chash

import (
	"crypto/md5"
	"fmt"
	"github.com/mayankz/gods/maps/hashmap"
	"github.com/mayankz/gods/maps/treemap"
	"github.com/mayankz/gods/utils"
)

type Node struct {
	Name   string
	Weight uint
}

type HashRing struct {
	Circle *treemap.Map
	Nodes  *hashmap.Map
	Weight uint
}

func (h *HashRing) PrintStruct() {
	fmt.Printf("+%v", h)
}

func New() *HashRing {
	var h HashRing
	h.Circle = treemap.NewWith(utils.UInt32Comparator)
	h.Nodes = hashmap.New()
	return &h
}

func (h *HashRing) GetNode(key string) (node string, found bool) {
	bkeyi, ok := h.Circle.GetCeiling(getKetamaKey(key))
	bkey := bkeyi.(uint32)
	if ok {
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
	h.Weight += weight
	h.updateCircle(node, true)
}

func (h *HashRing) RemoveNode(node string) {
	weight, ok := h.Nodes.Get(node)
	if !ok {
		return
	}

	h.updateCircle(node, false)
	h.Nodes.Remove(node)
	h.Weight -= weight.(uint)
}

func (h *HashRing) UpdateNode(node string, weight uint) {
	h.RemoveNode(node)
	h.AddNode(node, weight)
}

func (h *HashRing) getBucketsCount(node string) uint {
	weight, ok := h.Nodes.Get(node)
	if !ok {
		return 0
	}
	return 40 * weight.(uint)
}

func (h *HashRing) updateCircle(node string, add bool) {
	bucketCount := h.getBucketsCount(node)
	for i := (uint)(0); i < bucketCount; i++ {
		bucket := fmt.Sprintf("%s-%d", node, i)
		keys := getKetamaKeys(bucket)
		for _, key := range keys {
			if add {
				h.Circle.Put(key, node)
			} else {
				h.Circle.Remove(key)
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
