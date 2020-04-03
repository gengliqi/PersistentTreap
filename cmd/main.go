package main

import (
	persistent_treap "github.com/gengliqi/persistent_treap/persistent_treap"
)

type Key uint
type Value uint

func (a Key) Equals(b persistent_treap.Equitable) bool {
	return a == b.(Key)
}

func (a Key) Less(b persistent_treap.Sortable) bool {
	return a < b.(Key)
}

func (a Value) Equals(b persistent_treap.Equitable) bool {
	return a == b.(Value)
}

const Count = 100

func main() {
	treaps := make([]*persistent_treap.PersistentTreap, int(Count+1))
	treaps[0] = persistent_treap.NewPersistentTreap()

	for i := 1; i <= Count; i++ {
		treaps[i] = treaps[i-1].Insert(Key(i), Value(i))
	}

	for i := 1; i <= Count; i++ {
		// we can find the key value inserted in previous version
		for j := 1; j <= i; j++ {
			val, ok := treaps[i].GetValue(Key(j))
			if !ok || !val.(Value).Equals(Value(j)) {
				panic("something wrong")
			}
		}
		// the latter insert key value cannot find in this version
		for j := i + 1; j <= Count; j++ {
			_, ok := treaps[i].GetValue(Key(j))
			if ok {
				panic("something wrong")
			}
		}
	}
}
