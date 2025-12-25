package repository_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/what-writers-like/backend/internal/domain"
	"github.com/what-writers-like/backend/internal/repository/gorm"
	"github.com/what-writers-like/backend/internal/testutils"
)

func TestWorkRepository_Create(t *testing.T) {
	db, cleanup := testutils.SetupTestDB(t)
	defer cleanup()

	writerRepo := gorm.NewWriterRepository(db)
	workRepo := gorm.NewWorkRepository(db)

	writer := domain.NewWriter(1, "Jane Austen", 1775, nil, nil)
	require.NoError(t, writerRepo.Create(writer))

	work := domain.NewWork(1, "Pride and Prejudice", 1)
	err := workRepo.Create(work)
	assert.NoError(t, err)
}

func TestWorkRepository_GetByID(t *testing.T) {
	db, cleanup := testutils.SetupTestDB(t)
	defer cleanup()

	writerRepo := gorm.NewWriterRepository(db)
	workRepo := gorm.NewWorkRepository(db)

	writer := domain.NewWriter(1, "Jane Austen", 1775, nil, nil)
	require.NoError(t, writerRepo.Create(writer))

	work := domain.NewWork(1, "Pride and Prejudice", 1)
	require.NoError(t, workRepo.Create(work))

	found, err := workRepo.GetByID(1)
	require.NoError(t, err)
	assert.Equal(t, work.ID(), found.ID())
	assert.Equal(t, work.Title(), found.Title())
	assert.Equal(t, work.AuthorID(), found.AuthorID())
}

func TestWorkRepository_GetByAuthorID(t *testing.T) {
	db, cleanup := testutils.SetupTestDB(t)
	defer cleanup()

	writerRepo := gorm.NewWriterRepository(db)
	workRepo := gorm.NewWorkRepository(db)

	writer := domain.NewWriter(1, "Jane Austen", 1775, nil, nil)
	require.NoError(t, writerRepo.Create(writer))

	work1 := domain.NewWork(1, "Pride and Prejudice", 1)
	work2 := domain.NewWork(2, "Sense and Sensibility", 1)
	require.NoError(t, workRepo.Create(work1))
	require.NoError(t, workRepo.Create(work2))

	works, err := workRepo.GetByAuthorID(1)
	require.NoError(t, err)
	assert.Len(t, works, 2)
}

func TestWorkRepository_List(t *testing.T) {
	db, cleanup := testutils.SetupTestDB(t)
	defer cleanup()

	writerRepo := gorm.NewWriterRepository(db)
	workRepo := gorm.NewWorkRepository(db)

	writer := domain.NewWriter(1, "Jane Austen", 1775, nil, nil)
	require.NoError(t, writerRepo.Create(writer))

	work1 := domain.NewWork(1, "Pride and Prejudice", 1)
	work2 := domain.NewWork(2, "Sense and Sensibility", 1)
	require.NoError(t, workRepo.Create(work1))
	require.NoError(t, workRepo.Create(work2))

	works, err := workRepo.List(10, 0)
	require.NoError(t, err)
	assert.Len(t, works, 2)
}

func TestWorkRepository_Update(t *testing.T) {
	db, cleanup := testutils.SetupTestDB(t)
	defer cleanup()

	writerRepo := gorm.NewWriterRepository(db)
	workRepo := gorm.NewWorkRepository(db)

	writer := domain.NewWriter(1, "Jane Austen", 1775, nil, nil)
	require.NoError(t, writerRepo.Create(writer))

	work := domain.NewWork(1, "Pride and Prejudice", 1)
	require.NoError(t, workRepo.Create(work))

	updated := domain.NewWork(1, "Pride and Prejudice (Revised)", 1)
	err := workRepo.Update(updated)
	require.NoError(t, err)

	found, err := workRepo.GetByID(1)
	require.NoError(t, err)
	assert.Equal(t, "Pride and Prejudice (Revised)", found.Title())
}

func TestWorkRepository_Delete(t *testing.T) {
	db, cleanup := testutils.SetupTestDB(t)
	defer cleanup()

	writerRepo := gorm.NewWriterRepository(db)
	workRepo := gorm.NewWorkRepository(db)

	writer := domain.NewWriter(1, "Jane Austen", 1775, nil, nil)
	require.NoError(t, writerRepo.Create(writer))

	work := domain.NewWork(1, "Pride and Prejudice", 1)
	require.NoError(t, workRepo.Create(work))

	err := workRepo.Delete(1)
	require.NoError(t, err)

	_, err = workRepo.GetByID(1)
	assert.Error(t, err)
}
