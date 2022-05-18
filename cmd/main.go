// Copyright 2022 gengliqi.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
	treaps := make([]persistent_treap.PersistentTreap, int(Count+1))
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
