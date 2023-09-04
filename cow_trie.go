package main

// implementation of a copy-on-write trie data structure

type COWNode struct {
	Value string
	Next  map[rune]*COWNode
}

func NewCOWNode() *COWNode {
	return &COWNode{
		Value: "",
		Next:  make(map[rune]*COWNode),
	}
}

// Clone node value and all the children next pointers map
func (node *COWNode) Clone() *COWNode {
	clone := &COWNode{
		Value: node.Value,
		Next:  make(map[rune]*COWNode),
	}

	for k, v := range node.Next {
		clone.Next[k] = v
	}

	return clone
}

// Sets key - value pair recursively creating a new cloned node along the way
func (node *COWNode) Set(key string, value string) (*COWNode, error) {

	if key == "" {
		node.Value = value
		return node, nil
	}

	current := node

	runeStr := []rune(key)
	r := runeStr[0]

	currentClone := current.Clone()
	next, ok := currentClone.Next[r]
	if !ok {
		next = NewCOWNode()
	}

	next, err := next.Set(string(runeStr[1:]), value)
	if err != nil {
		return nil, err
	}

	currentClone.Next[r] = next

	return currentClone, nil
}

// Get method
func (node *COWNode) Get(key string) (string, error) {
	current := node
	for _, r := range key {
		next, ok := current.Next[r]
		if !ok {
			return "", ErrKeyNotFound
		}
		current = next
	}
	if current.Value != "" {
		return current.Value, nil
	}
	return "", ErrKeyNotFound
}

func (node *COWNode) Delete(key string) (*COWNode, error) {
	return deleteCOWNode(node, key, nil, 0)
}

// Delete function. Works recursively and deletes nodes if they become empty
// An empty node contains no children map pointers and empty value
func deleteCOWNode(node *COWNode, key string, parent *COWNode, parentKey rune) (*COWNode, error) {

	if key == "" {
		if parent == nil {
			return nil, ErrParentDelNotAllowed
		}
		nodeClone := node.Clone()
		nodeClone.Value = ""

		if len(nodeClone.Next) == 0 {
			delete(parent.Next, parentKey)
			return nil, nil
		}

		return nodeClone, nil
	}

	runeStr := []rune(key)
	r := runeStr[0]

	nodeClone := node.Clone()
	next, ok := node.Next[r]
	if !ok {
		return nil, ErrKeyNotFound
	}

	next, err := deleteCOWNode(next, string(runeStr[1:]), nodeClone, r)
	if err != nil {
		return nil, err
	}

	nodeClone.Next[r] = next
	if next == nil {
		delete(nodeClone.Next, r)
	}

	if parent != nil && len(nodeClone.Next) == 0 && nodeClone.Value == "" {
		delete(parent.Next, parentKey)
		return nil, nil
	}

	return nodeClone, nil
}
