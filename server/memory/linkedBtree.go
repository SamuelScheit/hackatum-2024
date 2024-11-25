package memory

import (
	"sync"

	"github.com/puzpuzpuz/xsync/v3"
	"github.com/tidwall/btree"
)

type LinkedBtree struct {
	root              btree.Map[int32, []int32]
	cacheGreaterEqual *xsync.MapOf[int32, *BitArray]
	cacheLessThan     *xsync.MapOf[int32, *BitArray]
	cacheLessEqual    *xsync.MapOf[int32, *BitArray]
	Size              int
	mutex             sync.Mutex
}

func NewLinkedBtree() *LinkedBtree {

	return &LinkedBtree{
		root: btree.Map[int32, []int32]{},
		cacheGreaterEqual: xsync.NewMapOf[int32, *BitArray](
			xsync.WithGrowOnly(),
		),
		cacheLessThan: xsync.NewMapOf[int32, *BitArray](
			xsync.WithGrowOnly(),
		),
		cacheLessEqual: xsync.NewMapOf[int32, *BitArray](
			xsync.WithGrowOnly(),
		),
		Size:  0,
		mutex: sync.Mutex{},
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

	t.cacheGreaterEqual.Range(func(k int32, v *BitArray) bool {
		if k >= key {
			v.SetBit(int(value))
		}
		return true
	})

	t.cacheLessThan.Range(func(k int32, v *BitArray) bool {
		if k < key {
			v.SetBit(int(value))
		}
		return true
	})

	t.cacheLessEqual.Range(func(k int32, v *BitArray) bool {
		if k <= key {
			v.SetBit(int(value))
		}
		return true
	})

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

	if cache, exists := t.cacheGreaterEqual.Load(key); exists {
		return cache.Copy()
	}

	resBitArray := NewBitArray(t.Size)
	t.cacheGreaterEqual.Store(key, resBitArray)

	t.GreaterThanEqual(key, func(key int32, iids []int32) {
		for _, iid := range iids {
			resBitArray.SetBit(int(iid))
		}
	})

	return resBitArray
}

func (t *LinkedBtree) BitArrayLessThan(key int32) *BitArray {
	if cache, exists := t.cacheLessThan.Load(key); exists {
		return cache.Copy()
	}

	resBitArray := NewBitArray(t.Size)

	t.cacheLessThan.Store(key, resBitArray)

	t.LessThan(key, func(key2 int32, iids []int32) {
		for _, iid := range iids {
			resBitArray.SetBit(int(iid))
		}
	})
	return resBitArray
}

func (t *LinkedBtree) BitArrayLessEqual(key int32) *BitArray {
	if cache, exists := t.cacheLessEqual.Load(key); exists {
		return cache.Copy()
	}

	resBitArray := NewBitArray(t.Size)

	t.cacheLessEqual.Store(key, resBitArray)

	t.LessThanEqual(key, func(key int32, iids []int32) {
		for _, iid := range iids {
			resBitArray.SetBit(int(iid))
		}
	})

	return resBitArray
}
