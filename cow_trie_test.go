package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var cowTests = []struct {
	key   string
	value string
}{
	{"key1", "value1"},
	{"key2", "value2"},
	{"key3", "value3"},
	{"key4", "value3"},
	{"key5", "value3"},
	{"a", "1"},
	{"ab", "2"},
	{"abc", "3"},
	{"abcd", "4"},
	{"abcde", "5"},
	{"abcdef", "6"},
	{"abcdefg", "7"},
	{"abcdefgh", "8"},
}

func TestCOW(t *testing.T) {

	t.Run("get non existing key", func(t *testing.T) {
		key := "non existing key"
		root := NewCOWNode()

		_, err := root.Get(key)
		require.Error(t, err)
	})

	t.Run("set and get values", func(t *testing.T) {

		root := NewCOWNode()
		var err error
		for _, tt := range cowTests {
			root, err = root.Set(tt.key, tt.value)
			require.NoError(t, err)

			value, err := root.Get(tt.key)
			require.NoError(t, err)
			require.Equal(t, tt.value, value)
		}

		for _, tt := range cowTests {
			value, err := root.Get(tt.key)
			require.NoError(t, err)
			require.Equal(t, tt.value, value)
		}
	})

	t.Run("set and get reverse values", func(t *testing.T) {

		root := NewCOWNode()
		var err error
		for i := len(cowTests) - 1; i >= 0; i-- {
			tt := cowTests[i]
			root, err = root.Set(tt.key, tt.value)
			require.NoError(t, err)

			value, err := root.Get(tt.key)
			require.NoError(t, err)
			require.Equal(t, tt.value, value)
		}

		for _, tt := range cowTests {
			value, err := root.Get(tt.key)
			require.NoError(t, err)
			require.Equal(t, tt.value, value)
		}

	})

	t.Run("delete non existing key", func(t *testing.T) {
		key := "non existing key"
		root := NewCOWNode()

		_, err := root.Delete(key)
		require.Error(t, err)
	})

	t.Run("cannot delete root", func(t *testing.T) {
		root := NewCOWNode()

		_, err := root.Delete("")
		require.Error(t, err)
	})

	t.Run("set and delete values", func(t *testing.T) {
		root := NewCOWNode()
		var err error
		for _, tt := range cowTests {
			root, err = root.Set(tt.key, tt.value)
			require.NoError(t, err)
		}

		for _, tt := range cowTests {
			root, err = root.Delete(tt.key)
			require.NoError(t, err)

			_, err = root.Get(tt.key)
			require.Error(t, err)
		}

		for _, tt := range cowTests {
			_, err = root.Get(tt.key)
			require.Error(t, err)
		}
	})

	t.Run("delete and get in order", func(t *testing.T) {
		root := NewCOWNode()
		var err error

		// Set values
		for _, tt := range cowTests {
			root, err = root.Set(tt.key, tt.value)
			require.NoError(t, err)
		}

		// Delete values one-by-one and check
		for i := 0; i < len(cowTests); i++ {
			tt := cowTests[i]

			root, err = root.Delete(tt.key)
			require.NoError(t, err)

			for j := i + 1; j < len(cowTests); j++ {
				value, err := root.Get(cowTests[j].key)
				require.NoError(t, err)
				require.Equal(t, cowTests[j].value, value)
			}

			_, err = root.Get(tt.key)
			require.Error(t, err)
		}
	})

	t.Run("delete and get in reverse order", func(t *testing.T) {
		root := NewCOWNode()
		var err error

		// Set values
		for _, tt := range cowTests {
			root, err = root.Set(tt.key, tt.value)
			require.NoError(t, err)
		}

		// Delete values in reverse order and check
		for i := len(cowTests) - 1; i >= 0; i-- {
			tt := cowTests[i]

			root, err = root.Delete(tt.key)
			require.NoError(t, err)

			for j := i - 1; j >= 0; j-- {
				value, err := root.Get(cowTests[j].key)
				require.NoError(t, err)
				require.Equal(t, cowTests[j].value, value)
			}

			_, err = root.Get(tt.key)
			require.Error(t, err)
		}
	})
}
