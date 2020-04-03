package persistent_treap

import (
	"math/rand"
	"sort"
	"testing"
)

type Key uint
type Value uint

func (a Key) Equals(b Equitable) bool {
	return a == b.(Key)
}

func (a Key) Less(b Sortable) bool {
	return a < b.(Key)
}

func (a Value) Equals(b Equitable) bool {
	return a == b.(Value)
}

type KV struct {
	k Key
	v Value
}

type KVSlice []KV

func (s KVSlice) Len() int           { return len(s) }
func (s KVSlice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s KVSlice) Less(i, j int) bool { return s[i].k.Less(s[j].k) }

func checkIsSame(t *testing.T, treap PersistentTreap, mp map[Key]Value) {
	var kvSlice KVSlice
	for k, v := range mp {
		kvSlice = append(kvSlice, KV{Key(k), Value(v)})
	}
	sort.Sort(kvSlice)

	allKV := treap.GetAllKeyValue()
	for i := 0; i < len(mp); i++ {
		if (*allKV[i].Key).(Key) != kvSlice[i].k || (*allKV[i].Value).(Value) != kvSlice[i].v {
			t.Fatalf("not equal! index:%v, tree key:%v, map key:%v, tree value:%v, map value:%v", i, (*allKV[i].Key).(Key), kvSlice[i].k, (*allKV[i].Value).(Value), kvSlice[i].v)
		}
	}
}

func randomInsertRemoveGetCount(t *testing.T, Count int) {
	keyRange := uint32(100000)

	mp := make(map[Key]Value)
	treaps := make([]PersistentTreap, Count+1)
	treaps[0] = NewPersistentTreap()
	var maps []map[Key]Value
	var mapPos []int

	t.Logf("start")
	for i := 0; i < Count; i++ {
		pos := rand.Uint32() % 100
		if pos <= 60 {
			k := rand.Uint32() % keyRange
			v := rand.Uint32()
			treaps[i+1] = treaps[i].Insert(Key(k), Value(v))
			mp[Key(k)] = Value(v)
		} else if pos <= 80 {
			k := rand.Uint32() % keyRange
			treaps[i+1] = treaps[i].Remove(Key(k))
			delete(mp, Key(k))
		} else {
			k := rand.Uint32() % keyRange
			val, ok := treaps[i].GetValue(Key(k))
			val2, ok2 := mp[Key(k)]
			if ok != ok2 {
				t.Fatalf("get value, key:%v, tree ok:%v, map ok:%v, tree val:%v, map val:%v", k, ok, ok2, val.(Value), val2)
			} else if ok {
				if val.(Value) != val2 {
					t.Fatalf("value not equal, key:%v, tree val:%v, map val:%v", k, val, val2)
				}
			}
			treaps[i+1] = treaps[i]
		}
		if i != 0 && i%100 == 3 {
			copyMap := make(map[Key]Value)
			for k, v := range mp {
				copyMap[k] = v
			}
			maps = append(maps, copyMap)
			mapPos = append(mapPos, i)
		}
	}
	t.Logf("end")
	t.Logf("map len:%v tree size:%v", len(mp), treaps[Count].Size())
	if uint32(len(mp)) != treaps[Count].Size() {
		t.Fatalf("size not equal!")
	}
	for i := 0; i < len(mapPos); i++ {
		checkIsSame(t, treaps[mapPos[i]+1], maps[i])
	}
	t.Logf("check ok")
}

func TestRandomInsertRemoveGet10000(t *testing.T) {
	randomInsertRemoveGetCount(t, 10000)
}

func TestRandomInsertRemoveGet100000(t *testing.T) {
	randomInsertRemoveGetCount(t, 100000)
}

func TestIsSameTreap(t *testing.T) {
	Count := 100000
	keyRange := uint32(100000)

	mp := make(map[Key]Value)
	treap1 := NewPersistentTreap()
	treap2 := NewPersistentTreap()
	treap3 := NewPersistentTreap()

	for i := 0; i < Count; i++ {
		k := rand.Uint32() % keyRange
		v := rand.Uint32()
		mp[Key(k)] = Value(v)
		treap1 = treap1.Insert(Key(k), Value(v))
	}
	for k, v := range mp {
		treap2 = treap2.Insert(Key(k), Value(v))
	}
	isEqual := true
	for k, v := range mp {
		if rand.Uint32()%100 < 99 {
			isEqual = false
			treap3 = treap3.Insert(Key(k), Value(v))
		}
	}
	if !IsSameTreap(treap1, treap2) {
		t.Fatalf("treap1 and treap2 must be same!")
	}
	if isEqual != IsSameTreap(treap1, treap3) {
		t.Fatalf("treap1 and treap3 isEqual:%v fail", isEqual)
	}
	if isEqual != IsSameTreap(treap2, treap3) {
		t.Fatalf("treap2 and treap3 isEqual:%v fail", isEqual)
	}
}
