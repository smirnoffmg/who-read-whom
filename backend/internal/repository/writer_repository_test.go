package repository_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/what-writers-like/backend/internal/domain"
	"github.com/what-writers-like/backend/internal/repository/gorm"
	"github.com/what-writers-like/backend/internal/testutils"
)

func TestWriterRepository_Create(t *testing.T) {
	db, cleanup := testutils.SetupTestDB(t)
	defer cleanup()

	repo := gorm.NewWriterRepository(db)

	writer := domain.NewWriter(1, "Jane Austen", 1775, nil, nil)
	err := repo.Create(writer)
	assert.NoError(t, err)
}

func TestWriterRepository_GetByID(t *testing.T) {
	db, cleanup := testutils.SetupTestDB(t)
	defer cleanup()

	repo := gorm.NewWriterRepository(db)

	writer := domain.NewWriter(1, "Jane Austen", 1775, nil, nil)
	err := repo.Create(writer)
	require.NoError(t, err)

	found, err := repo.GetByID(1)
	require.NoError(t, err)
	assert.Equal(t, writer.ID(), found.ID())
	assert.Equal(t, writer.Name(), found.Name())
	assert.Equal(t, writer.BirthYear(), found.BirthYear())
}

func TestWriterRepository_List(t *testing.T) {
	db, cleanup := testutils.SetupTestDB(t)
	defer cleanup()

	repo := gorm.NewWriterRepository(db)

	writer1 := domain.NewWriter(1, "Jane Austen", 1775, nil, nil)
	writer2 := domain.NewWriter(2, "Charles Dickens", 1812, nil, nil)

	require.NoError(t, repo.Create(writer1))
	require.NoError(t, repo.Create(writer2))

	writers, err := repo.List(10, 0)
	require.NoError(t, err)
	assert.Len(t, writers, 2)
}

func TestWriterRepository_Update(t *testing.T) {
	db, cleanup := testutils.SetupTestDB(t)
	defer cleanup()

	repo := gorm.NewWriterRepository(db)

	writer := domain.NewWriter(1, "Jane Austen", 1775, nil, nil)
	require.NoError(t, repo.Create(writer))

	bio := "English novelist"
	updated := domain.NewWriter(1, "Jane Austen", 1775, nil, &bio)
	err := repo.Update(updated)
	require.NoError(t, err)

	found, err := repo.GetByID(1)
	require.NoError(t, err)
	assert.Equal(t, bio, *found.Bio())
}

func TestWriterRepository_Delete(t *testing.T) {
	db, cleanup := testutils.SetupTestDB(t)
	defer cleanup()

	repo := gorm.NewWriterRepository(db)

	writer := domain.NewWriter(1, "Jane Austen", 1775, nil, nil)
	require.NoError(t, repo.Create(writer))

	err := repo.Delete(1)
	require.NoError(t, err)

	_, err = repo.GetByID(1)
	assert.Error(t, err)
}
