package adapters

import (
	"encoding/json"
	"fmt"
	"os"
	domain "spotify_migration/entities"
	"spotify_migration/entities/data"
)

const (
	pendingState int = iota
	migratedState
	deletedState
)

type state struct {
	State   int    `json:"state"`
	Address string `json:"address"`
}

func migrated(state int) bool {
	return state == migratedState
}

type MigrationState struct {
	filename string
	items    map[string]state
}

func NewMigrationState(collectionID string) domain.IMigrationStateRepository {
	return &MigrationState{
		filename: fmt.Sprintf("%s_migration_state.json", collectionID),
		items:    make(map[string]state),
	}
}

func (m *MigrationState) GetPendingItems() map[string]string {
	pendingItems := make(map[string]string)

	for id, state := range m.items {
		if !migrated(state.State) {
			pendingItems[id] = state.Address
		}
	}
	return pendingItems
}

func (m *MigrationState) UpdateItemToMigrated(itemID string) {
	if m.items != nil {
		currentState, exists := m.items[itemID]

		if !exists {
			return
		}
		m.items[itemID] = state{
			State:   migratedState,
			Address: currentState.Address,
		}
	}
}

// func (m *MigrationState) UpdateItemToMigrated(item *data.Music) {
// 	if m.items != nil {
// 		currentState, exists := m.items[domain.ID(item)]

// 		if !exists {
// 			return
// 		}
// 		m.items[domain.ID(item)] = state{
// 			state:   migratedState,
// 			address: currentState.address,
// 		}
// 	}
// }

func (m *MigrationState) AddItem(item *data.Music, address string) {
	if m.items != nil {

		itemID := domain.ID(item)

		if _, exists := m.items[itemID]; exists {
			return
		}

		m.items[itemID] = state{
			State:   pendingState,
			Address: address,
		}
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
