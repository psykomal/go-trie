package main

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

var storeTests = []struct {
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

var stores = []struct {
	name string
}{
	{name: "MapStore"},
	{name: "TrieStore"},
}

func StoreTest(storeType string) func(t *testing.T) {

	return func(t *testing.T) {

		t.Run("empty key set, get, delete", func(t *testing.T) {
			store := GetStore(storeType)
			key := ""
			value := "val"

			err := store.Set(key, value)
			require.ErrorIs(t, err, ErrKeyIsEmpty)

			_, err = store.Get(key)
			require.ErrorIs(t, err, ErrKeyIsEmpty)

			err = store.Delete(key)
			require.ErrorIs(t, err, ErrKeyIsEmpty)
		})

		t.Run("get non existent key", func(t *testing.T) {
			store := GetStore(storeType)
			key := "non key"

			_, err := store.Get(key)
			require.ErrorIs(t, err, ErrKeyNotFound)
		})

		t.Run("Set and Get", func(t *testing.T) {
			store := GetStore(storeType)

			for _, tt := range storeTests {
				err := store.Set(tt.key, tt.value)
				require.NoError(t, err)

				value, err := store.Get(tt.key)
				require.NoError(t, err)
				require.Equal(t, tt.value, value)
			}

			for _, tt := range storeTests {
				value, err := store.Get(tt.key)
				require.NoError(t, err)
				require.Equal(t, tt.value, value)
			}
		})

		t.Run("Set duplicate key serially", func(t *testing.T) {
			store := GetStore(storeType)
			key := "key"
			valuePrefix := "val"
			var lastValue string

			for i := 0; i < 100; i++ {
				value := valuePrefix + strconv.Itoa(i)
				lastValue = value
				err := store.Set(key, value)
				require.NoError(t, err)
			}

			lastValue, err := store.Get(key)
			require.NoError(t, err)
			require.Equal(t, lastValue, lastValue)
		})

		t.Run("Set single key parallelly", func(t *testing.T) {
			store := GetStore(storeType)
			key := "key"
			valuePrefix := "val"
			var lastValue string
			var wg sync.WaitGroup
			var mu sync.Mutex

			for i := 0; i < 100; i++ {
				wg.Add(1)
				go func(i int) {
					defer wg.Done()
					mu.Lock()
					defer mu.Unlock()

					value := valuePrefix + strconv.Itoa(i)
					err := store.Set(key, value)
					require.NoError(t, err)
					lastValue = value
				}(i)
			}

			wg.Wait()

			value, err := store.Get(key)
			require.NoError(t, err)
			require.Equal(t, lastValue, value)
		})

		t.Run("Set, get and delete keys parallelly", func(t *testing.T) {
			store := GetStore(storeType)
			var wg sync.WaitGroup

			for i := 0; i < len(storeTests); i++ {
				wg.Add(1)

				go func(i int) {
					defer wg.Done()
					tt := storeTests[i]

					err := store.Set(tt.key, tt.value)
					require.NoError(t, err)
				}(i)
			}

			wg.Wait()

			for i := 0; i < len(storeTests); i++ {
				wg.Add(1)

				go func(i int) {
					defer wg.Done()
					tt := storeTests[i]

					value, err := store.Get(tt.key)
					require.NoError(t, err)
					require.Equal(t, value, tt.value)
				}(i)
			}

			wg.Wait()

			for i := 0; i < len(storeTests); i++ {
				wg.Add(1)

				go func(i int) {
					defer wg.Done()
					tt := storeTests[i]

					err := store.Delete(tt.key)
					require.NoError(t, err)
				}(i)
			}

			wg.Wait()
		})

		t.Run("Run all Set, get and delete parallelly", func(t *testing.T) {
			store := GetStore(storeType)
			var wg sync.WaitGroup
			getKeys := make(chan string)
			delKeys := make(chan string)

			for i := 0; i < len(storeTests); i++ {
				wg.Add(1)

				go func(i int) {
					defer wg.Done()
					tt := storeTests[i]

					err := store.Set(tt.key, tt.value)
					require.NoError(t, err)

					getKeys <- tt.key
				}(i)
			}

			for i := 0; i < len(storeTests); i++ {
				wg.Add(1)

				go func() {
					defer wg.Done()

					key := <-getKeys

					_, err := store.Get(key)
					require.NoError(t, err)

					delKeys <- key
				}()
			}

			for i := 0; i < len(storeTests); i++ {
				wg.Add(1)

				go func() {
					defer wg.Done()

					key := <-delKeys

					err := store.Delete(key)
					require.NoError(t, err)
				}()
			}

			wg.Wait()
			close(getKeys)
			close(delKeys)
		})
	}
}

func TestStores(t *testing.T) {

	for _, tt := range stores {
		storeType := tt.name
		t.Run(storeType, func(t *testing.T) {
			StoreTest(storeType)(t)
		})
	}
}

func BenchmarkStores(b *testing.B) {

	for _, store := range stores {

		b.Run(store.name, func(b *testing.B) {
			store := GetStore(store.name)

			b.ResetTimer()

			b.RunParallel(func(pb *testing.PB) {
				var key string
				var value string
				for pb.Next() {
					key = strconv.Itoa(rand.Int())
					value = strconv.Itoa(rand.Int())
					err := store.Set(key, value)
					require.NoError(b, err)

					_, _ = store.Get(key)

					_ = store.Delete(key)
				}
			})
		})
	}
}
