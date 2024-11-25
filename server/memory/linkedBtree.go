package memory

import (
	"github.com/tidwall/btree"
)

type LinkedBtree struct {
	root              btree.Map[int32, []int32]
	cacheGreaterEqual map[int32]*BitArray
	cacheLessThan     map[int32]*BitArray
	cacheLessEqual    map[int32]*BitArray
	Size              int
}

func NewLinkedBtree() *LinkedBtree {

	return &LinkedBtree{
		root:              btree.Map[int32, []int32]{},
		cacheGreaterEqual: map[int32]*BitArray{},
		cacheLessThan:     map[int32]*BitArray{},
		cacheLessEqual:    map[int32]*BitArray{},
		Size:              0,
	}
}

func (t *LinkedBtree) Add(key, value int32) {
	keys, exists := t.root.Get(key)
	if !exists {
		keys = []int32{}
	}

	keys = append(keys, value)

	t.root.Set(key, keys)
	t.Size++

	for k, v := range t.cacheGreaterEqual {
		if k >= key {
			v.SetBit(int(value))
		}
	}

	for k, v := range t.cacheLessThan {
		if k < key {
			v.SetBit(int(value))
		}
	}

	for k, v := range t.cacheLessEqual {
		if k <= key {
			v.SetBit(int(value))
		}
	}
}

func (t *LinkedBtree) GreaterThan(key int32, fn func(int32, []int32)) {
	t.root.Ascend(key, func(a int32, b []int32) bool {
		if key == a {
			return true
		}

		fn(a, b)
		return true
	})
}

func (t *LinkedBtree) GreaterThanEqual(key int32, fn func(int32, []int32)) {
	t.root.Ascend(key, func(a int32, b []int32) bool {

		fn(a, b)
		return true
	})
}

func (t *LinkedBtree) LessThan(key int32, fn func(int32, []int32)) {
	t.root.Descend(key, func(a int32, b []int32) bool {
		if key == a {
			return true
		}
		fn(a, b)
		return true
	})
}

func (t *LinkedBtree) LessThanEqual(key int32, fn func(int32, []int32)) {

	t.root.Descend(key, func(a int32, b []int32) bool {
		fn(a, b)
		return true
	})
}

func (t *LinkedBtree) BitArrayGreaterEqual(key int32) *BitArray {
	if cache, exists := t.cacheGreaterEqual[key]; exists {
		return cache.Copy()
	}

	resBitArray := NewBitArray(t.Size)

	t.cacheGreaterEqual[key] = resBitArray

	t.GreaterThanEqual(key, func(key int32, iids []int32) {
		for _, iid := range iids {
			resBitArray.SetBit(int(iid))
		}
	})

	return resBitArray
}

func (t *LinkedBtree) BitArrayLessThan(key int32) *BitArray {
	if cache, exists := t.cacheLessThan[key]; exists {
		return cache.Copy()
	}

	resBitArray := NewBitArray(t.Size)

	t.cacheLessThan[key] = resBitArray

	t.LessThan(key, func(key2 int32, iids []int32) {
		for _, iid := range iids {
			resBitArray.SetBit(int(iid))
		}
	})
	return resBitArray
}

func (t *LinkedBtree) BitArrayLessEqual(key int32) *BitArray {
	if cache, exists := t.cacheLessEqual[key]; exists {
		return cache.Copy()
	}

	resBitArray := NewBitArray(t.Size)

	t.cacheLessEqual[key] = resBitArray

	t.LessThanEqual(key, func(key int32, iids []int32) {
		for _, iid := range iids {
			resBitArray.SetBit(int(iid))
		}
	})

	return resBitArray
}
