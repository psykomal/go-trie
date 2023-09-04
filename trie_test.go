package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var tests = []struct {
	key   string
	value string
}{
	{"key1", "value1"},
	{"key2", "value2"},
	{"key3", "value3"},
	{"a", "1"},
	{"ab", "2"},
	{"abc", "3"},
}

var deleteTests = []struct {
	key   string
	value string
}{
	{"key2", "value2"},
	{"a", "1"},
}

var afterDeleteTests = []struct {
	key   string
	value string
}{
	{"key1", "value1"},
	{"key3", "value3"},
	{"ab", "2"},
	{"abc", "3"},
}

func setUp(t *testing.T) *Node {
	root := NewNode()
	for _, test := range tests {
		err := root.Put(test.key, test.value)
		require.NoError(t, err)
	}
	return root
}

func TestPutAndGet(t *testing.T) {

	root := setUp(t)

	for _, test := range tests {

		value, ok := root.Get(test.key)
		require.True(t, ok)
		require.Equal(t, test.value, value)
	}
}

func TestDelete(t *testing.T) {
	root := setUp(t)

	for _, test := range deleteTests {
		err := root.Delete(test.key)
		require.NoError(t, err)
	}

	for _, test := range deleteTests {
		_, ok := root.Get(test.key)
		require.False(t, ok)
	}

	for _, test := range afterDeleteTests {
		value, ok := root.Get(test.key)
		require.True(t, ok)
		require.Equal(t, test.value, value)
	}
}
