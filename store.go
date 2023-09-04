package main

type Store interface {
	Set(key string, value string) error
	Get(key string) (string, error)
	Delete(key string) error
}

func GetStore(name string) Store {
	switch name {
	case "MapStore":
		return NewMapStore()
	case "TrieStore":
		return NewTrieStore()
	default:
		panic("unknown store type")
	}
}
