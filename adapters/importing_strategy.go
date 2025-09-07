package adapters

func NewYoutubeMemsetUpdater() *YoutubeMemsetUpdater {
	return &YoutubeMemsetUpdater{}
}

type YoutubeMemsetUpdater struct {
}

func (u *YoutubeMemsetUpdater) UpdateItems(collectionID string, itemIDs []string) error {
	err := u.deleteAll(collectionID)
	if err != nil {
		return err
	}

	err = u.insertAll(collectionID, itemIDs)
	if err != nil {
		return err
	}

	return nil
}

func (u *YoutubeMemsetUpdater) deleteAll(collectionID string) error {
	// Simulate deleting all items from the playlist
	if collectionID == "" {
		return nil
	}
	return nil
}

func (u *YoutubeMemsetUpdater) insertAll(collectionID string, itemIDs []string) error {
	if collectionID == "" || len(itemIDs) == 0 {
		return nil
	}
	for _, itemID := range itemIDs {
		err := u.addItemToPlaylist(collectionID, itemID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (u *YoutubeMemsetUpdater) addItemToPlaylist(collectionID string, itemID string) error {
	// Simulate adding an item to the playlist
	if collectionID == "" || itemID == "" {
		return nil
	}
	return nil
}
