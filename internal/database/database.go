package database

import (
	"sync"
)

type Database struct {
	mu      sync.RWMutex
	strings map[string]string
	lists   map[string][]string
}

func (db *Database) Unlock() {
	panic("unimplemented")
}

func (db *Database) Lock() {
	panic("unimplemented")
}

func (db *Database) SetStrings(strings map[string]string) {
	panic("unimplemented")
}

func New() *Database {
	return &Database{
		strings: make(map[string]string),
		lists:   make(map[string][]string),
	}
}

// Métodos para strings
func (db *Database) Set(key, value string) {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.strings[key] = value
}

func (db *Database) Get(key string) (string, bool) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	val, exists := db.strings[key]
	return val, exists
}

// Métodos para listas
func (db *Database) LPush(key, value string) int {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.lists[key] = append([]string{value}, db.lists[key]...)
	return len(db.lists[key])
}

func (db *Database) SetLists(l map[string][]string) {
	db.lists = l
}

func (db *Database) GetStrings() map[string]string {
	return db.strings
}

func (db *Database) GetLists() map[string][]string {
	return db.lists
}
