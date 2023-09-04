package main

import "fmt"

type Node struct {
	Value string
	Next  map[rune]*Node
}

func NewNode() *Node {
	return &Node{
		Value: "",
		Next:  make(map[rune]*Node),
	}
}

func (root *Node) Put(key string, value string) error {
	current := root
	for _, r := range key {
		next, ok := current.Next[r]
		if !ok {
			next = NewNode()
			current.Next[r] = next
		}
		current = next
	}
	current.Value = value

	return nil
}

func (root *Node) Get(key string) (string, bool) {
	current := root
	for _, r := range key {
		next, ok := current.Next[r]
		if !ok {
			return "", false
		}
		current = next
	}
	if current.Value != "" {
		return current.Value, true
	}
	return "", false
}

func (root *Node) Delete(key string) error {
	current := root
	var prev *Node
	for _, r := range key {
		prev = current
		next, ok := current.Next[r]
		if !ok {
			return fmt.Errorf(key + " not found")
		}
		current = next
	}

	current.Value = ""

	if len(current.Next) == 0 && prev != nil {
		delete(prev.Next, []rune(key)[0])
	}

	return nil
}
