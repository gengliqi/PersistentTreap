package main

import (
	"fmt"

	"github.com/gengliqi/PersistentTreap"
)

type key int
type value int

func (a key) Equitable(b key) bool {
	return a == b
}

func (a key) Less(b key) bool {
	return a < b
}

func (a value) Equitable(b value) bool {
	return a == b
}

func main() {
	fmt.Printf("haha");
	treap := persistenttreap.NewPersistentTreap()
	treap.Insert(key(5), value(10))
	println("5:%u", treap.GetValue(key(5)))
}