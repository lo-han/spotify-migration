package data

type IMigrationStateRepository interface {
	GetPendingItems() map[string]string
	GetMigratedItems() map[string]string
	UpdateItemToMigrated(itemID string)
	AddItem(item *Music, address string)
	Read() (bool, error)
	Save() error
}
