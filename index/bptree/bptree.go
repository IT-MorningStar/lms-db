package bptree

import "lms-db/engine/storage"

type BPTree struct {
	root *Node
	file storage.FileAccess
}

type Node struct {
	id   int64
	leaf bool
}

func NewRootNode() *Node {
	return nil
}

func NewBPTree(root *Node) *BPTree {
	return &BPTree{
		root: root,
	}
}

func (bp *BPTree) GetByNodeId(nodeId int64) Node {
	return Node{}
}
