package memory

import (
	"github.com/puzpuzpuz/xsync/v3"
	"github.com/tidwall/btree"
)

type LinkedBtree struct {
	root              btree.Map[int32, []int32]
	cacheGreaterEqual *xsync.MapOf[int32, *BitArray]
	cacheLessThan     *xsync.MapOf[int32, *BitArray]
	cacheLessEqual    *xsync.MapOf[int32, *BitArray]
	Size              int
}

func hashInt32(k int32, _ uint64) uint64 {
	return uint64(k)
}

func NewLinkedBtree() *LinkedBtree {

	return &LinkedBtree{
		root: btree.Map[int32, []int32]{},
		cacheGreaterEqual: xsync.NewMapOfWithHasher[int32, *BitArray](
			hashInt32,
			xsync.WithGrowOnly(),
		),
		cacheLessThan: xsync.NewMapOfWithHasher[int32, *BitArray](
			hashInt32,
			xsync.WithGrowOnly(),
		),
		cacheLessEqual: xsync.NewMapOfWithHasher[int32, *BitArray](
			hashInt32,
			xsync.WithGrowOnly(),
		),
		Size: 0,
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

	t.cacheGreaterEqual.Delete(key)
	t.cacheLessThan.Delete(key)
	t.cacheLessEqual.Delete(key)

	// t.cacheGreaterEqual.Range(func(k int32, v *BitArray) bool {
	// 	if k >= key {
	// 		v.SetBit(int(value))
	// 	}
	// 	return true
	// })

	// t.cacheLessThan.Range(func(k int32, v *BitArray) bool {
	// 	if k < key {
	// 		v.SetBit(int(value))
	// 	}
	// 	return true
	// })

	// t.cacheLessEqual.Range(func(k int32, v *BitArray) bool {
	// 	if k <= key {
	// 		v.SetBit(int(value))
	// 	}
	// 	return true
	// })

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
	value, loaded := t.cacheGreaterEqual.LoadOrCompute(key, func() *BitArray {
		return NewBitArray(t.Size)
	})

	if !loaded {
		t.GreaterThanEqual(key, func(key2 int32, iids []int32) {
			for _, iid := range iids {
				value.SetBit(int(iid))
			}
		})
	}

	return value.Copy()
}

func (t *LinkedBtree) BitArrayLessThan(key int32) *BitArray {
	value, loaded := t.cacheLessThan.LoadOrCompute(key, func() *BitArray {
		return NewBitArray(t.Size)
	})

	if !loaded {
		t.LessThan(key, func(key int32, iids []int32) {
			for _, iid := range iids {
				value.SetBit(int(iid))
			}
		})
	}

	return value.Copy()
}

func (t *LinkedBtree) BitArrayLessEqual(key int32) *BitArray {
	value, loaded := t.cacheLessEqual.LoadOrCompute(key, func() *BitArray {
		return NewBitArray(t.Size)
	})

	if !loaded {
		t.LessThanEqual(key, func(key int32, iids []int32) {
			for _, iid := range iids {
				value.SetBit(int(iid))
			}
		})
	}

	return value.Copy()
}
