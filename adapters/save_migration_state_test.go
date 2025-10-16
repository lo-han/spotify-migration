package adapters

import (
	"fmt"
	"os"
	"reflect"
	"testing"
)

func newMigrationStateTest(collectionID string) *MigrationState {
	return &MigrationState{
		filename: fmt.Sprintf("%s_migration_state.json", collectionID),
		items:    make(map[string]bool),
	}
}

func Test_MigrationIsStarting(t *testing.T) {
	t.Run("MigrationIsStarting", func(t *testing.T) {
		ms := newMigrationStateTest("test-migration-starting")
		ms.items["song1"] = false
		ms.items["song2"] = false
		ms.items["song3"] = false
		ms.items["song4"] = false

		if err := ms.Save(); err != nil {
			t.Fatalf("Failed to save migration state: %v", err)
		}

		testFilename := ms.filename

		defer func() {
			if err := os.Remove(testFilename); err != nil && !os.IsNotExist(err) {
				t.Logf("Warning: could not clean up test file %s: %v", testFilename, err)
			}
		}()

		msRead := newMigrationStateTest("test-migration-starting")

		exists, err := msRead.Read()
		if err != nil {
			t.Fatalf("Failed to read migration state: %v", err)
		}

		if !exists {
			t.Error("Expected file to exist after save")
		}

		pendingItems := msRead.GetPendingItems()

		expectedPending := []string{"song1", "song2", "song3", "song4"}

		if !reflect.DeepEqual(pendingItems, expectedPending) {
			t.Errorf("Expected pending items [song3 song4], got %v", msRead.GetPendingItems())
		}
	})
}

func Test_MigrationIsInProcess(t *testing.T) {
	t.Run("Test_MigrationIsInProcess", func(t *testing.T) {
		ms := newMigrationStateTest("test-save-read")
		ms.items["song1"] = true
		ms.items["song2"] = true
		ms.items["song3"] = false
		ms.items["song4"] = false
		ms.items["song5"] = false

		ms.UpdateItemToMigrated("song3")

		if err := ms.Save(); err != nil {
			t.Fatalf("Failed to save migration state: %v", err)
		}

		testFilename := ms.filename

		defer func() {
			if err := os.Remove(testFilename); err != nil && !os.IsNotExist(err) {
				t.Logf("Warning: could not clean up test file %s: %v", testFilename, err)
			}
		}()

		msRead := newMigrationStateTest("test-save-read")

		exists, err := msRead.Read()
		if err != nil {
			t.Fatalf("Failed to read migration state: %v", err)
		}

		if !exists {
			t.Error("Expected file to exist after save")
		}

		pendingItems := msRead.GetPendingItems()

		expectedPending := []string{"song4", "song5"}

		if !reflect.DeepEqual(pendingItems, expectedPending) {
			t.Errorf("Expected pending items [song3 song4], got %v", msRead.GetPendingItems())
		}
	})
}

func Test_MigrationHasFinished(t *testing.T) {
	t.Run("Test_MigrationHasFinished", func(t *testing.T) {

		ms := newMigrationStateTest("test-migration-finished")
		ms.items["song1"] = true
		ms.items["song2"] = true
		ms.items["song3"] = true
		ms.items["song4"] = true

		if err := ms.Save(); err != nil {
			t.Fatalf("Failed to save migration state: %v", err)
		}

		testFilename := ms.filename

		defer func() {
			if err := os.Remove(testFilename); err != nil && !os.IsNotExist(err) {
				t.Logf("Warning: could not clean up test file %s: %v", testFilename, err)
			}
		}()

		msRead := newMigrationStateTest("test-migration-finished")

		exists, err := msRead.Read()
		if err != nil {
			t.Fatalf("Failed to read migration state: %v", err)
		}

		if !exists {
			t.Error("Expected file to exist after save")
		}

		if !reflect.DeepEqual(ms.items, msRead.items) {
			t.Errorf("Expected items %v, got %v", ms.items, msRead.items)
		}

		pendingItems := msRead.GetPendingItems()
		if len(pendingItems) != 0 {
			t.Errorf("Expected no pending items for finished migration, got %v", pendingItems)
		}
	})
}
