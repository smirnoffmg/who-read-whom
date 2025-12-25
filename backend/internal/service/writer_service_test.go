package service_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/what-writers-like/backend/internal/domain"
	"github.com/what-writers-like/backend/internal/repository/gorm"
	"github.com/what-writers-like/backend/internal/service"
	"github.com/what-writers-like/backend/internal/testutils"
)

func TestWriterService_CreateWriter(t *testing.T) {
	t.Parallel()
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		db, cleanup := testutils.SetupTestDB(t)
		defer cleanup()

		writerRepo := gorm.NewWriterRepository(db)
		workRepo := gorm.NewWorkRepository(db)
		svc := service.NewWriterService(writerRepo, workRepo)

		writer, err := svc.CreateWriter("Jane Austen", 1775, nil, nil)
		require.NoError(t, err)
		assert.Equal(t, "Jane Austen", writer.Name())
		assert.Equal(t, 1775, writer.BirthYear())
	})

	t.Run("empty name", func(t *testing.T) {
		t.Parallel()
		db, cleanup := testutils.SetupTestDB(t)
		defer cleanup()

		writerRepo := gorm.NewWriterRepository(db)
		workRepo := gorm.NewWorkRepository(db)
		svc := service.NewWriterService(writerRepo, workRepo)

		_, err := svc.CreateWriter("", 1775, nil, nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "name is required")
	})

	t.Run("invalid birth year", func(t *testing.T) {
		t.Parallel()
		db, cleanup := testutils.SetupTestDB(t)
		defer cleanup()

		writerRepo := gorm.NewWriterRepository(db)
		workRepo := gorm.NewWorkRepository(db)
		svc := service.NewWriterService(writerRepo, workRepo)

		_, err := svc.CreateWriter("Jane Austen", 0, nil, nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "birth year must be positive")
	})

	t.Run("increments id correctly", func(t *testing.T) {
		t.Parallel()
		db, cleanup := testutils.SetupTestDB(t)
		defer cleanup()

		writerRepo := gorm.NewWriterRepository(db)
		workRepo := gorm.NewWorkRepository(db)
		svc := service.NewWriterService(writerRepo, workRepo)

		// Create writers with specific IDs
		writer1 := domain.NewWriter(1, "Writer 1", 1800, nil, nil)
		writer3 := domain.NewWriter(3, "Writer 3", 1800, nil, nil)
		require.NoError(t, writerRepo.Create(writer1))
		require.NoError(t, writerRepo.Create(writer3))

		writer, err := svc.CreateWriter("New Writer", 1900, nil, nil)
		require.NoError(t, err)
		assert.Equal(t, uint64(4), writer.ID())
	})
}

func TestWriterService_GetWriter(t *testing.T) {
	t.Parallel()
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		db, cleanup := testutils.SetupTestDB(t)
		defer cleanup()

		writerRepo := gorm.NewWriterRepository(db)
		workRepo := gorm.NewWorkRepository(db)
		svc := service.NewWriterService(writerRepo, workRepo)

		expectedWriter := domain.NewWriter(1, "Jane Austen", 1775, nil, nil)
		require.NoError(t, writerRepo.Create(expectedWriter))

		writer, err := svc.GetWriter(1)
		require.NoError(t, err)
		assert.Equal(t, expectedWriter.ID(), writer.ID())
		assert.Equal(t, expectedWriter.Name(), writer.Name())
	})

	t.Run("not found", func(t *testing.T) {
		t.Parallel()
		db, cleanup := testutils.SetupTestDB(t)
		defer cleanup()

		writerRepo := gorm.NewWriterRepository(db)
		workRepo := gorm.NewWorkRepository(db)
		svc := service.NewWriterService(writerRepo, workRepo)

		_, err := svc.GetWriter(999)
		require.Error(t, err)
	})
}

func TestWriterService_ListWriters(t *testing.T) {
	t.Parallel()
	db, cleanup := testutils.SetupTestDB(t)
	defer cleanup()

	writerRepo := gorm.NewWriterRepository(db)
	workRepo := gorm.NewWorkRepository(db)
	svc := service.NewWriterService(writerRepo, workRepo)

	writer1 := domain.NewWriter(1, "Jane Austen", 1775, nil, nil)
	writer2 := domain.NewWriter(2, "Charles Dickens", 1812, nil, nil)
	require.NoError(t, writerRepo.Create(writer1))
	require.NoError(t, writerRepo.Create(writer2))

	writers, err := svc.ListWriters(10, 0)
	require.NoError(t, err)
	assert.Len(t, writers, 2)
}

func TestWriterService_UpdateWriter(t *testing.T) {
	t.Parallel()
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		db, cleanup := testutils.SetupTestDB(t)
		defer cleanup()

		writerRepo := gorm.NewWriterRepository(db)
		workRepo := gorm.NewWorkRepository(db)
		svc := service.NewWriterService(writerRepo, workRepo)

		writer := domain.NewWriter(1, "Jane Austen", 1775, nil, nil)
		require.NoError(t, writerRepo.Create(writer))

		bio := "English novelist"
		err := svc.UpdateWriter(1, "Jane Austen", 1775, nil, &bio)
		require.NoError(t, err)

		updated, err := writerRepo.GetByID(1)
		require.NoError(t, err)
		assert.NotNil(t, updated.Bio())
		assert.Equal(t, bio, *updated.Bio())
	})

	t.Run("empty name", func(t *testing.T) {
		t.Parallel()
		db, cleanup := testutils.SetupTestDB(t)
		defer cleanup()

		writerRepo := gorm.NewWriterRepository(db)
		workRepo := gorm.NewWorkRepository(db)
		svc := service.NewWriterService(writerRepo, workRepo)

		err := svc.UpdateWriter(1, "", 1775, nil, nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "name is required")
	})

	t.Run("invalid birth year", func(t *testing.T) {
		t.Parallel()
		db, cleanup := testutils.SetupTestDB(t)
		defer cleanup()

		writerRepo := gorm.NewWriterRepository(db)
		workRepo := gorm.NewWorkRepository(db)
		svc := service.NewWriterService(writerRepo, workRepo)

		err := svc.UpdateWriter(1, "Jane Austen", 0, nil, nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "birth year must be positive")
	})
}

func TestWriterService_DeleteWriter(t *testing.T) {
	t.Parallel()
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		db, cleanup := testutils.SetupTestDB(t)
		defer cleanup()

		writerRepo := gorm.NewWriterRepository(db)
		workRepo := gorm.NewWorkRepository(db)
		svc := service.NewWriterService(writerRepo, workRepo)

		writer := domain.NewWriter(1, "Jane Austen", 1775, nil, nil)
		require.NoError(t, writerRepo.Create(writer))

		err := svc.DeleteWriter(1)
		require.NoError(t, err)

		_, err = writerRepo.GetByID(1)
		require.Error(t, err)
	})

	t.Run("cannot delete writer with works", func(t *testing.T) {
		t.Parallel()
		db, cleanup := testutils.SetupTestDB(t)
		defer cleanup()

		writerRepo := gorm.NewWriterRepository(db)
		workRepo := gorm.NewWorkRepository(db)
		svc := service.NewWriterService(writerRepo, workRepo)

		writer := domain.NewWriter(1, "Jane Austen", 1775, nil, nil)
		require.NoError(t, writerRepo.Create(writer))
		work := domain.NewWork(1, "Pride and Prejudice", 1)
		require.NoError(t, workRepo.Create(work))

		err := svc.DeleteWriter(1)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "cannot delete writer with existing works")
	})
}
