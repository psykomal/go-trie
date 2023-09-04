package main

import "errors"

// trie errors

var (
	ErrKeyNotFound         = errors.New("key not found")
	ErrParentDelNotAllowed = errors.New("deleting parent node not allowed")
)

// store errors

var (
	ErrKeyIsEmpty = errors.New("key is empty")
)
