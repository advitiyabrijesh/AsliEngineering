package main

import (
	"fmt"
	"sort"
)

type LSMTree struct {
	memTable  map[string]string
	diskTable map[string]string
}

func NewLSMTree() *LSMTree {
	return &LSMTree{
		memTable:  make(map[string]string),
		diskTable: make(map[string]string),
	}
}

func (tree *LSMTree) Put(key, value string) {
	tree.memTable[key] = value
}

func (tree *LSMTree) Get(key string) (string, bool, string) {
	value, ok := tree.memTable[key]
	if ok {
		return value, true, "Key found in mem table."
	}
	value, ok = tree.diskTable[key]
	var existenceMsg string
	if ok {
		existenceMsg = "Key found in disk table."
	} else {
		existenceMsg = "Key not found."
	}
	return value, ok, existenceMsg
}

func (tree *LSMTree) Compaction() {
	keys := make([]string, 0, len(tree.memTable))
	for key := range tree.memTable {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		value := tree.memTable[key]
		tree.diskTable[key] = value
		delete(tree.memTable, key)
	}
	fmt.Println("MemTable flushed to disk, current size of memTable: ", len(tree.memTable))
}

func main() {
	tree := NewLSMTree()
	tree.Put("k1", "v1")
	tree.Put("k2", "v2")

	value, ok, existenceMsg := tree.Get("k1")
	if ok {
		fmt.Println(value, existenceMsg)
	}

	tree.Compaction()

	value, ok, existenceMsg = tree.Get("k2")
	if ok {
		fmt.Println(value, existenceMsg)
	}

}
