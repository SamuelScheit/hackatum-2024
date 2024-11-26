package memory

import (
	"sync"

	"github.com/puzpuzpuz/xsync/v3"
	"github.com/tidwall/btree"
)

type LinkedBtree struct {
	root              btree.Map[int32, []int32]
	cGreaterEqual     map[int32]*BitArray
	cLessThan         map[int32]*BitArray
	cLessEqual        map[int32]*BitArray
	cacheGreaterEqual *xsync.MapOf[int32, *BitArray]
	cacheLessThan     *xsync.MapOf[int32, *BitArray]
	cacheLessEqual    *xsync.MapOf[int32, *BitArray]
	Size              int
	mutex             sync.Mutex
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
		cGreaterEqual: map[int32]*BitArray{},
		cLessThan:     map[int32]*BitArray{},
		cLessEqual:    map[int32]*BitArray{},
		Size:          0,
		mutex:         sync.Mutex{},
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

	// t.cacheGreaterEqual.Delete(key)
	// t.cacheLessThan.Delete(key)
	// t.cacheLessEqual.Delete(key)

	// delete(t.cGreaterEqual, key)
	// delete(t.cLessThan, key)
	// delete(t.cLessEqual, key)

	// for k, v := range t.cGreaterEqual {
	// 	if k >= key {
	// 		v.SetBit(int(value))
	// 	}
	// }

	// for k, v := range t.cLessThan {
	// 	if k < key {
	// 		v.SetBit(int(value))
	// 	}
	// }

	// for k, v := range t.cLessEqual {
	// 	if k <= key {
	// 		v.SetBit(int(value))
	// 	}
	// }

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
	value, _ := t.cacheGreaterEqual.LoadOrCompute(key, func() *BitArray {
		value := NewBitArray(t.Size)
		t.GreaterThanEqual(key, func(key2 int32, iids []int32) {
			for _, iid := range iids {
				value.SetBit(int(iid))
			}
		})

		return value
	})

	return value
}

func (t *LinkedBtree) BitArrayLessThan(key int32) *BitArray {
	value, _ := t.cacheLessThan.LoadOrCompute(key, func() *BitArray {
		value := NewBitArray(t.Size)

		t.LessThan(key, func(key int32, iids []int32) {
			for _, iid := range iids {
				value.SetBit(int(iid))
			}
		})
		return value
	})

	return value
}

func (t *LinkedBtree) BitArrayLessEqual(key int32) *BitArray {
	value, _ := t.cacheLessEqual.LoadOrCompute(key, func() *BitArray {
		value := NewBitArray(t.Size)

		t.LessThanEqual(key, func(key int32, iids []int32) {
			for _, iid := range iids {
				value.SetBit(int(iid))
			}
		})
		return value
	})

	return value
}

// func (t *LinkedBtree) BitArrayGreaterEqual(key int32) *BitArray {

// 	value := NewBitArray(t.Size)
// 	t.GreaterThanEqual(key, func(key2 int32, iids []int32) {
// 		for _, iid := range iids {
// 			value.SetBit(int(iid))
// 		}
// 	})

// 	return value
// }

// func (t *LinkedBtree) BitArrayLessThan(key int32) *BitArray {
// 	value := NewBitArray(t.Size)

// 	t.LessThan(key, func(key int32, iids []int32) {
// 		for _, iid := range iids {
// 			value.SetBit(int(iid))
// 		}
// 	})

// 	return value
// }

// func (t *LinkedBtree) BitArrayLessEqual(key int32) *BitArray {
// 	value := NewBitArray(t.Size)

// 	t.LessThanEqual(key, func(key int32, iids []int32) {
// 		for _, iid := range iids {
// 			value.SetBit(int(iid))
// 		}
// 	})

// 	return value
// }

// func (t *LinkedBtree) BitArrayGreaterEqual(key int32) *BitArray {
// 	if value, exists := t.cGreaterEqual[key]; exists {
// 		return value
// 	}

// 	value := NewBitArray(t.Size)
// 	t.GreaterThanEqual(key, func(key2 int32, iids []int32) {
// 		for _, iid := range iids {
// 			value.SetBit(int(iid))
// 		}
// 	})

// 	t.cGreaterEqual[key] = value

// 	return value
// }

// func (t *LinkedBtree) BitArrayLessThan(key int32) *BitArray {
// 	if value, exists := t.cLessThan[key]; exists {
// 		return value
// 	}

// 	value := NewBitArray(t.Size)

// 	t.LessThan(key, func(key int32, iids []int32) {
// 		for _, iid := range iids {
// 			value.SetBit(int(iid))
// 		}
// 	})

// 	t.cLessThan[key] = value

// 	return value
// }

// func (t *LinkedBtree) BitArrayLessEqual(key int32) *BitArray {
// 	if value, exists := t.cLessEqual[key]; exists {
// 		return value
// 	}

// 	value := NewBitArray(t.Size)

// 	t.LessThanEqual(key, func(key int32, iids []int32) {
// 		for _, iid := range iids {
// 			value.SetBit(int(iid))
// 		}
// 	})

// 	t.cLessEqual[key] = value

// 	return value
// }
