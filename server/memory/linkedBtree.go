package memory

import (
	"github.com/tidwall/btree"
)

type LinkedBtree struct {
	root btree.Map[int32, []int32]
	Size int
}

func NewLinkedBtree() *LinkedBtree {
	return &LinkedBtree{
		root: btree.Map[int32, []int32]{},
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
	t.LessThan(key, func(key2 int32, iids []int32) {
		for _, iid := range iids {
			resBitArray.SetBit(int(iid))
		}
	})
}

func (t *LinkedBtree) BitArrayLessEqual(key int32, resBitArray *BitArray) {
	t.LessThanEqual(key, func(key int32, iids []int32) {
		for _, iid := range iids {
			resBitArray.SetBit(int(iid))
		}
	})
}
