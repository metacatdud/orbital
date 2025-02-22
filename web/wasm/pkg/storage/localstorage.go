package storage

import (
	"encoding/json"
	"fmt"
	"syscall/js"
)

type LocalStorage struct {
	store js.Value
}

func (storage *LocalStorage) Get(key string, value any) error {
	entry := storage.store.Call("getItem", key)
	if !entry.Truthy() {
		return fmt.Errorf("%w:[%s]", ErrNotFound, key)
	}

	err := json.Unmarshal([]byte(entry.String()), value)
	if err != nil {
		return err
	}

	return nil
}

func (storage *LocalStorage) Set(key string, value any) error {

	b, err := json.Marshal(value)
	if err != nil {
		return err
	}

	storage.store.Call("setItem", key, string(b))

	return nil
}

func (storage *LocalStorage) Del(key string) error {
	storage.store.Call("removeItem", key)
	return nil
}

func (storage *LocalStorage) Exist(key string) bool {
	entry := storage.store.Call("getItem", key)
	if !entry.Truthy() {
		return false
	}
	return true
}

func NewLocalStorage() *LocalStorage {
	return &LocalStorage{
		store: js.Global().Get("localStorage"),
	}
}
