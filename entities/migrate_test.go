package entities

import (
	"context"
	"errors"
	"testing"

	"spotify_migration/entities/data"
	"spotify_migration/mocks"

	"github.com/stretchr/testify/assert"
)

func TestMigration_Migrate(t *testing.T) {
	t.Run("should migrate resource successfully", func(t *testing.T) {
		ctx := context.Background()
		resourceName := "test_playlist"

		extractorMock := mocks.NewIExtractorUsecase(t)
		importerMock := mocks.NewIImporterUsecase(t)

		migration := NewMigration(extractorMock, importerMock)

		collection := &data.Collection{}

		extractorMock.On("Extract", ctx, resourceName).Return(collection, nil)
		importerMock.On("Import", ctx, collection).Return(true, nil)

		success, err := migration.Migrate(ctx, resourceName)

		assert.NoError(t, err)
		assert.True(t, success)
	})

	t.Run("should return false if export failed", func(t *testing.T) {
		ctx := context.Background()
		resourceName := "test_playlist"

		extractorMock := mocks.NewIExtractorUsecase(t)
		importerMock := mocks.NewIImporterUsecase(t)

		migration := NewMigration(extractorMock, importerMock)

		collection := &data.Collection{}

		extractorMock.On("Extract", ctx, resourceName).Return(collection, nil)
		importerMock.On("Import", ctx, collection).Return(false, nil)

		success, err := migration.Migrate(ctx, resourceName)

		assert.NoError(t, err)
		assert.False(t, success)
	})

	t.Run("should return false if import returned false", func(t *testing.T) {
		ctx := context.Background()
		resourceName := "test_playlist"

		extractorMock := mocks.NewIExtractorUsecase(t)
		importerMock := mocks.NewIImporterUsecase(t)

		migration := NewMigration(extractorMock, importerMock)

		collection := &data.Collection{}

		extractorMock.On("Extract", ctx, resourceName).Return(collection, nil)
		importerMock.On("Import", ctx, collection).Return(false, nil)

		success, err := migration.Migrate(ctx, resourceName)

		assert.NoError(t, err)
		assert.False(t, success)
	})

	t.Run("should return false if import failed", func(t *testing.T) {
		ctx := context.Background()
		resourceName := "test_playlist"

		extractorMock := mocks.NewIExtractorUsecase(t)
		importerMock := mocks.NewIImporterUsecase(t)

		migration := NewMigration(extractorMock, importerMock)

		collection := &data.Collection{}

		expectedError := errors.New("import error")

		extractorMock.On("Extract", ctx, resourceName).Return(collection, nil)
		importerMock.On("Import", ctx, collection).Return(true, expectedError)

		succeeded, err := migration.Migrate(ctx, resourceName)

		assert.ErrorIs(t, err, expectedError)
		assert.False(t, succeeded)
	})

}
