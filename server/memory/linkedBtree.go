package memory

import (
	"github.com/tidwall/btree"
)

type LinkedBtree struct {
	root    btree.Map[int32, []int32]
	hashmap map[int32][]int32
	Size    int
}

func NewLinkedBtree() *LinkedBtree {
	return &LinkedBtree{
		root:    btree.Map[int32, []int32]{},
		hashmap: make(map[int32][]int32),
	}
}

func (t *LinkedBtree) Add(key, value int32) {
	if list, ok := t.hashmap[key]; ok {
		t.hashmap[key] = append(list, value)
	} else {
		t.hashmap[key] = []int32{value}
	}

	t.root.Set(key, t.hashmap[key])
	t.Size++
}

func (t *LinkedBtree) Get(key int32) []int32 {
	if list, ok := t.hashmap[key]; ok {
		return list
	}
	return []int32{}
}

func (t *LinkedBtree) GreaterThan(key int32, fn func(int32, []int32)) {
	t.root.Ascend(key, func(a int32, b []int32) bool {
		fn(a, b)
		return true
	})
}

func (t *LinkedBtree) GreaterThanEqual(key int32, fn func(int32, []int32)) {
	if list, ok := t.hashmap[key]; ok {
		fn(key, list)
	}

	t.GreaterThan(key, fn)
}

func (t *LinkedBtree) LessThan(key int32, fn func(int32, []int32)) {
	t.root.Descend(key, func(a int32, b []int32) bool {
		fn(a, b)
		return true
	})
}

func (t *LinkedBtree) LessThanEqual(key int32, fn func(int32, []int32)) {
	if list, ok := t.hashmap[key]; ok {
		fn(key, list)
	}

	t.LessThan(key, fn)
}

func (t *LinkedBtree) BitArrayGreaterEqual(key int32, resBitArray *BitArray) {
	t.GreaterThanEqual(key, func(key int32, iids []int32) {
		for _, iid := range iids {
			resBitArray.SetBit(int(iid))
		}
	})
}

func (t *LinkedBtree) BitArrayGreaterThan(key int32, resBitArray *BitArray) {
	t.GreaterThan(key, func(key int32, iids []int32) {
		for _, iid := range iids {
			resBitArray.SetBit(int(iid))
		}
	})
}

func (t *LinkedBtree) BitArrayLessThan(key int32, resBitArray *BitArray) {
	t.LessThan(key, func(key int32, iids []int32) {
		for _, iid := range iids {
			resBitArray.SetBit(int(iid))
		}
	})
}

func (t *LinkedBtree) BitArrayLessThanEqual(key int32, resBitArray *BitArray) {
	t.LessThanEqual(key, func(key int32, iids []int32) {
		for _, iid := range iids {
			resBitArray.SetBit(int(iid))
		}
	})
}
