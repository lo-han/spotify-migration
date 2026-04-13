package entities

import (
	"spotify_migration/entities/data"
	"testing"
)

func TestID(t *testing.T) {
	t.Run("Should return empty if music is nil", func(t *testing.T) {
		var music *data.Music = nil
		result := ID(music)
		if result != "" {
			t.Errorf("Expected empty string for nil music, got %q", result)
		}
	})

	t.Run("Standard music with all fields", func(t *testing.T) {
		music := &data.Music{
			Name:   "Bohemian Rhapsody",
			Artist: "Queen",
			Album:  PtrStr("A Night at the Opera"),
		}

		expected := "Bohemian Rhapsody_Queen_A Night at the Opera"
		result := ID(music)

		if result != expected {
			t.Errorf("Expected ID %q, got %q", expected, result)
		}
	})

	t.Run("Music with empty title", func(t *testing.T) {
		music := &data.Music{
			Name:   "",
			Artist: "The Beatles",
			Album:  PtrStr("Abbey Road"),
		}

		expected := "_The Beatles_Abbey Road"
		result := ID(music)

		if result != expected {
			t.Errorf("Expected ID %q, got %q", expected, result)
		}
	})

	t.Run("Music with empty artist", func(t *testing.T) {
		music := &data.Music{
			Name:   "Imagine",
			Artist: "",
			Album:  PtrStr("Imagine"),
		}

		expected := "Imagine__Imagine"
		result := ID(music)

		if result != expected {
			t.Errorf("Expected ID %q, got %q", expected, result)
		}
	})

	t.Run("Music with empty album", func(t *testing.T) {
		music := &data.Music{
			Name:   "Hotel California",
			Artist: "Eagles",
			Album:  PtrStr(""),
		}

		expected := "Hotel California_Eagles_"
		result := ID(music)

		if result != expected {
			t.Errorf("Expected ID %q, got %q", expected, result)
		}
	})

	t.Run("Music with all empty fields", func(t *testing.T) {
		music := &data.Music{
			Name:   "",
			Artist: "",
			Album:  PtrStr(""),
		}

		expected := "__"
		result := ID(music)

		if result != expected {
			t.Errorf("Expected ID %q, got %q", expected, result)
		}
	})

	t.Run("Music with special characters", func(t *testing.T) {
		music := &data.Music{
			Name:   "Don't Stop Me Now",
			Artist: "Queen & David Bowie",
			Album:  PtrStr("Jazz (Deluxe Edition)"),
		}

		expected := "Don't Stop Me Now_Queen & David Bowie_Jazz (Deluxe Edition)"
		result := ID(music)

		if result != expected {
			t.Errorf("Expected ID %q, got %q", expected, result)
		}
	})

	t.Run("Music with underscores in fields", func(t *testing.T) {
		music := &data.Music{
			Name:   "Song_With_Underscores",
			Artist: "Artist_Name",
			Album:  PtrStr("Album_Title"),
		}

		expected := "Song_With_Underscores_Artist_Name_Album_Title"
		result := ID(music)

		if result != expected {
			t.Errorf("Expected ID %q, got %q", expected, result)
		}
	})

	t.Run("Music with whitespace", func(t *testing.T) {
		music := &data.Music{
			Name:   "  Stairway to Heaven  ",
			Artist: "  Led Zeppelin  ",
			Album:  PtrStr("  Led Zeppelin IV  "),
		}

		expected := "  Stairway to Heaven  _  Led Zeppelin  _  Led Zeppelin IV  "
		result := ID(music)

		if result != expected {
			t.Errorf("Expected ID %q, got %q", expected, result)
		}
	})
}

func TestID_Consistency(t *testing.T) {
	t.Run("Same music produces same ID", func(t *testing.T) {
		music := &data.Music{
			Name:   "Test Song",
			Artist: "Test Artist",
			Album:  PtrStr("Test Album"),
		}

		id1 := ID(music)
		id2 := ID(music)

		if id1 != id2 {
			t.Errorf("Expected same ID for same music, got %q and %q", id1, id2)
		}
	})

	t.Run("Different music with same fields produces same ID", func(t *testing.T) {
		music1 := &data.Music{
			Name:   "Test Song",
			Artist: "Test Artist",
			Album:  PtrStr("Test Album"),
		}

		music2 := &data.Music{
			Name:   "Test Song",
			Artist: "Test Artist",
			Album:  PtrStr("Test Album"),
		}

		id1 := ID(music1)
		id2 := ID(music2)

		if id1 != id2 {
			t.Errorf("Expected same ID for music with same fields, got %q and %q", id1, id2)
		}
	})

	t.Run("Different music produces different ID", func(t *testing.T) {
		music1 := &data.Music{
			Name:   "Test Song 1",
			Artist: "Test Artist",
			Album:  PtrStr("Test Album"),
		}

		music2 := &data.Music{
			Name:   "Test Song 2",
			Artist: "Test Artist",
			Album:  PtrStr("Test Album"),
		}

		id1 := ID(music1)
		id2 := ID(music2)

		if id1 == id2 {
			t.Errorf("Expected different IDs for different music, got same ID %q", id1)
		}
	})
}
