package database

import (
	"encoding/json"
	"errors"
	"fmt"
)

var ErrNoRecords = errors.New("no records found")

type InMemoryDB struct {
	storage map[string]map[int64][]byte
}

func NewInMemoryDB() InMemoryDB {
	return InMemoryDB{storage: make(map[string]map[int64][]byte)}
}

func (db InMemoryDB) Insert(tableName string, id int64, data interface{}) error {
	table, existsTable := db.storage[tableName]

	j, _ := json.Marshal(data)

	if !existsTable {
		db.storage[tableName] = make(map[int64][]byte)
		db.storage[tableName][id] = j
		return nil
	}

	if _, exists := table[id]; exists {
		return fmt.Errorf("there is already a data for id '%d'", id)
	}

	table[id] = j

	return nil
}

func (db InMemoryDB) Find(tableName string, id int64, target interface{}) error {
	table, existsTable := db.storage[tableName]

	if !existsTable {
		return ErrNoRecords
	}

	if v, exists := table[id]; exists {
		_ = json.Unmarshal(v, target)
		return nil
	}

	return ErrNoRecords
}

func (db InMemoryDB) Update(tableName string, id int64, data interface{}) error {
	table := db.storage[tableName]

	j, _ := json.Marshal(data)

	table[id] = j

	return nil
}
