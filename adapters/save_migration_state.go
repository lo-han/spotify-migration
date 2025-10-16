package adapters

import (
	"encoding/json"
	"fmt"
	"os"
	"spotify_migration/domain"
)

type MigrationState struct {
	filename string
	items    map[string]bool
}

func NewMigrationState(collectionID string) domain.IMigrationStateRepository {
	return &MigrationState{
		filename: fmt.Sprintf("%s_migration_state.json", collectionID),
		items:    make(map[string]bool),
	}
}

func (m *MigrationState) GetPendingItems() []string {
	pendingItems := []string{}

	for item, migrated := range m.items {
		if !migrated {
			pendingItems = append(pendingItems, item)
		}
	}
	return pendingItems
}

func (m *MigrationState) UpdateItemToMigrated(itemID string) {
	if m.items != nil {
		m.items[itemID] = true
	}
}

func (m *MigrationState) Read() (bool, error) {
	file, err := os.Open(m.filename)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&m.items); err != nil {
		return true, err
	}
	return true, nil
}

func (m *MigrationState) Save() error {
	file, err := os.Create(m.filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	return encoder.Encode(m.items)
}
