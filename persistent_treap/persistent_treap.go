package persistent_treap

import (
	"log"
	"math/rand"
)

type PersistentTreap struct {
	root *treapNode
}

func NewPersistentTreap() PersistentTreap {
	return PersistentTreap{root: nil}
}

func (tree PersistentTreap) Find(key Sortable) bool {
	return tree.root.find(key)
}

func (tree PersistentTreap) GetValue(key Sortable) (Equitable, bool) {
	return tree.root.getValue(key)
}

// Return a new tree whose data is equal to previous tree + key, value
func (tree PersistentTreap) Insert(key Sortable, value Equitable) PersistentTreap {
	newRoot := tree.root.insert(key, value)
	if newRoot == tree.root {
		return tree
	}
	return PersistentTreap{newRoot}
}

// Return a new tree whose data is equal to previous tree - key.
// If the key is not exist, the new tree's root is the same as previous one
func (tree PersistentTreap) Remove(key Sortable) PersistentTreap {
	newRoot := tree.root.remove(key)
	if newRoot == tree.root {
		return tree
	}
	return PersistentTreap{newRoot}
}

func (tree PersistentTreap) Size() uint32 {
	return tree.root.size()
}

func (tree PersistentTreap) Clear() PersistentTreap {
	return NewPersistentTreap()
}

type TreapKeyValueRef struct {
	Key   *Sortable
	Value *Equitable
}

func (u *PersistentTreap) GetAllKeyValue() []TreapKeyValueRef {
	var uKV []TreapKeyValueRef
	u.root.getAllKeyValue(&uKV)
	return uKV
}

func IsSameTreap(u, v PersistentTreap) bool {
	if u.Size() != v.Size() {
		return false
	}
	uKV := make([]TreapKeyValueRef, 0, u.Size())
	vKV := make([]TreapKeyValueRef, 0, u.Size())
	u.root.getAllKeyValue(&uKV)
	v.root.getAllKeyValue(&vKV)
	if len(uKV) != len(vKV) || len(uKV) != int(u.Size()) {
		log.Panicf("keyValue size not equal, len(u):%v, len(v):%v, size:%v", len(uKV), len(vKV), u.Size())
	}
	for i := 0; i < len(uKV); i++ {
		if !(*uKV[i].Key).Equals(*vKV[i].Key) || !(*uKV[i].Value).Equals(*vKV[i].Value) {
			return false
		}
	}
	return true
}

type Equitable interface {
	Equals(b Equitable) bool
}

type Sortable interface {
	Equitable
	Less(b Sortable) bool
}

type treapNode struct {
	left     *treapNode
	right    *treapNode
	key      Sortable
	val      Equitable
	priority uint32
	sz       uint32
}

func newTreapNode(key Sortable, value Equitable) *treapNode {
	return &treapNode{
		left:     nil,
		right:    nil,
		key:      key,
		val:      value,
		priority: rand.Uint32(),
		sz:       1,
	}
}

func (u *treapNode) updateSize() {
	if u != nil {
		u.sz = u.left.size() + u.right.size() + 1
	}
}

// If leftEqual is true, means the range is (-∞, key], (key, +∞)
// Otherwise the range is (-∞, key), [key, +∞)
func split(u *treapNode, key Sortable, leftEqual bool) (*treapNode, *treapNode) {
	if u == nil {
		return nil, nil
	}
	newRoot := *u
	if key.Less(newRoot.key) || (!leftEqual && key.Equals(newRoot.key)) {
		leftNode, rightNode := split(newRoot.left, key, leftEqual)
		newRoot.left = rightNode
		newRoot.updateSize()
		return leftNode, &newRoot
	} else {
		leftNode, rightNode := split(newRoot.right, key, leftEqual)
		newRoot.right = leftNode
		newRoot.updateSize()
		return &newRoot, rightNode
	}
}

func merge(u *treapNode, v *treapNode) *treapNode {
	if u == nil {
		return v
	} else if v == nil {
		return u
	}
	leftBigger := true
	if u.priority < v.priority {
		leftBigger = false
	} else if u.priority == v.priority {
		if rand.Uint32()%2 == 1 {
			leftBigger = false
		}
	}
	if leftBigger {
		u.right = merge(u.right, v)
		u.updateSize()
		return u
	} else {
		v.left = merge(u, v.left)
		v.updateSize()
		return v
	}
}

func (u *treapNode) find(key Sortable) bool {
	_, ret := u.getValue(key)
	return ret
}

func (u *treapNode) getValue(key Sortable) (Equitable, bool) {
	node := u
	for {
		if node == nil {
			return nil, false
		}
		if key.Equals(node.key) {
			return node.val, true
		} else if key.Less(node.key) {
			node = node.left
		} else {
			node = node.right
		}
	}
}

func (u *treapNode) insert(key Sortable, value Equitable) *treapNode {
	if u == nil {
		return newTreapNode(key, value)
	}
	oldValue, isFind := u.getValue(key)
	if isFind {
		if oldValue.Equals(value) {
			return u
		}
		root := *u
		node := &root
		for {
			if key.Equals(node.key) {
				node.val = value
				break
			} else if key.Less(node.key) {
				newLeft := *node.left
				node.left = &newLeft
				node = node.left
			} else {
				newRight := *node.right
				node.right = &newRight
				node = node.right
			}
		}
		return &root
	}
	newLeft, newRight := split(u, key, true)
	return merge(merge(newLeft, newTreapNode(key, value)), newRight)
}

func (u *treapNode) remove(key Sortable) *treapNode {
	if u == nil || !u.find(key) {
		return u
	}
	left1, right1 := split(u, key, false)
	_, right2 := split(right1, key, true)
	return merge(left1, right2)
}

func (u *treapNode) size() uint32 {
	if u == nil {
		return 0
	}
	return u.sz
}

func (u *treapNode) getAllKeyValue(s *[]TreapKeyValueRef) {
	if u == nil {
		return
	}
	u.left.getAllKeyValue(s)
	*s = append(*s, TreapKeyValueRef{&u.key, &u.val})
	u.right.getAllKeyValue(s)
}
