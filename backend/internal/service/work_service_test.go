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

func TestWorkService_CreateWork(t *testing.T) {
	t.Parallel()
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		db, cleanup := testutils.SetupTestDB(t)
		defer cleanup()

		workRepo := gorm.NewWorkRepository(db)
		writerRepo := gorm.NewWriterRepository(db)
		svc := service.NewWorkService(workRepo, writerRepo)

		writer := domain.NewWriter(1, "Jane Austen", 1775, nil, nil)
		require.NoError(t, writerRepo.Create(writer))

		work, err := svc.CreateWork("Pride and Prejudice", 1)
		require.NoError(t, err)
		assert.Equal(t, "Pride and Prejudice", work.Title())
		assert.Equal(t, uint64(1), work.AuthorID())
	})

	t.Run("empty title", func(t *testing.T) {
		t.Parallel()
		db, cleanup := testutils.SetupTestDB(t)
		defer cleanup()

		workRepo := gorm.NewWorkRepository(db)
		writerRepo := gorm.NewWriterRepository(db)
		svc := service.NewWorkService(workRepo, writerRepo)

		_, err := svc.CreateWork("", 1)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "title is required")
	})

	t.Run("author not found", func(t *testing.T) {
		t.Parallel()
		db, cleanup := testutils.SetupTestDB(t)
		defer cleanup()

		workRepo := gorm.NewWorkRepository(db)
		writerRepo := gorm.NewWriterRepository(db)
		svc := service.NewWorkService(workRepo, writerRepo)

		_, err := svc.CreateWork("Pride and Prejudice", 999)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "author not found")
	})
}

func TestWorkService_GetWork(t *testing.T) {
	t.Parallel()
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		db, cleanup := testutils.SetupTestDB(t)
		defer cleanup()

		workRepo := gorm.NewWorkRepository(db)
		writerRepo := gorm.NewWriterRepository(db)
		svc := service.NewWorkService(workRepo, writerRepo)

		expectedWork := domain.NewWork(1, "Pride and Prejudice", 1)
		require.NoError(t, workRepo.Create(expectedWork))

		work, err := svc.GetWork(1)
		require.NoError(t, err)
		assert.Equal(t, expectedWork.ID(), work.ID())
		assert.Equal(t, expectedWork.Title(), work.Title())
	})

	t.Run("not found", func(t *testing.T) {
		t.Parallel()
		db, cleanup := testutils.SetupTestDB(t)
		defer cleanup()

		workRepo := gorm.NewWorkRepository(db)
		writerRepo := gorm.NewWriterRepository(db)
		svc := service.NewWorkService(workRepo, writerRepo)

		_, err := svc.GetWork(999)
		require.Error(t, err)
	})
}

func TestWorkService_GetWorksByAuthor(t *testing.T) {
	t.Parallel()
	db, cleanup := testutils.SetupTestDB(t)
	defer cleanup()

	workRepo := gorm.NewWorkRepository(db)
	writerRepo := gorm.NewWriterRepository(db)
	svc := service.NewWorkService(workRepo, writerRepo)

	writer := domain.NewWriter(1, "Jane Austen", 1775, nil, nil)
	require.NoError(t, writerRepo.Create(writer))

	work1 := domain.NewWork(1, "Pride and Prejudice", 1)
	work2 := domain.NewWork(2, "Sense and Sensibility", 1)
	require.NoError(t, workRepo.Create(work1))
	require.NoError(t, workRepo.Create(work2))

	works, err := svc.GetWorksByAuthor(1)
	require.NoError(t, err)
	assert.Len(t, works, 2)
}

func TestWorkService_ListWorks(t *testing.T) {
	t.Parallel()
	db, cleanup := testutils.SetupTestDB(t)
	defer cleanup()

	workRepo := gorm.NewWorkRepository(db)
	writerRepo := gorm.NewWriterRepository(db)
	svc := service.NewWorkService(workRepo, writerRepo)

	work1 := domain.NewWork(1, "Pride and Prejudice", 1)
	work2 := domain.NewWork(2, "Sense and Sensibility", 1)
	require.NoError(t, workRepo.Create(work1))
	require.NoError(t, workRepo.Create(work2))

	works, err := svc.ListWorks(10, 0)
	require.NoError(t, err)
	assert.Len(t, works, 2)
}

func TestWorkService_UpdateWork(t *testing.T) {
	t.Parallel()
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		db, cleanup := testutils.SetupTestDB(t)
		defer cleanup()

		workRepo := gorm.NewWorkRepository(db)
		writerRepo := gorm.NewWriterRepository(db)
		svc := service.NewWorkService(workRepo, writerRepo)

		writer := domain.NewWriter(1, "Jane Austen", 1775, nil, nil)
		require.NoError(t, writerRepo.Create(writer))
		work := domain.NewWork(1, "Pride and Prejudice", 1)
		require.NoError(t, workRepo.Create(work))

		err := svc.UpdateWork(1, "Pride and Prejudice (Revised)", 1)
		require.NoError(t, err)

		updated, err := workRepo.GetByID(1)
		require.NoError(t, err)
		assert.Equal(t, "Pride and Prejudice (Revised)", updated.Title())
	})

	t.Run("empty title", func(t *testing.T) {
		t.Parallel()
		db, cleanup := testutils.SetupTestDB(t)
		defer cleanup()

		workRepo := gorm.NewWorkRepository(db)
		writerRepo := gorm.NewWriterRepository(db)
		svc := service.NewWorkService(workRepo, writerRepo)

		err := svc.UpdateWork(1, "", 1)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "title is required")
	})

	t.Run("work not found", func(t *testing.T) {
		t.Parallel()
		db, cleanup := testutils.SetupTestDB(t)
		defer cleanup()

		workRepo := gorm.NewWorkRepository(db)
		writerRepo := gorm.NewWriterRepository(db)
		svc := service.NewWorkService(workRepo, writerRepo)

		err := svc.UpdateWork(999, "Pride and Prejudice", 1)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "work not found")
	})

	t.Run("author not found", func(t *testing.T) {
		t.Parallel()
		db, cleanup := testutils.SetupTestDB(t)
		defer cleanup()

		workRepo := gorm.NewWorkRepository(db)
		writerRepo := gorm.NewWriterRepository(db)
		svc := service.NewWorkService(workRepo, writerRepo)

		writer := domain.NewWriter(1, "Jane Austen", 1775, nil, nil)
		require.NoError(t, writerRepo.Create(writer))
		work := domain.NewWork(1, "Pride and Prejudice", 1)
		require.NoError(t, workRepo.Create(work))

		err := svc.UpdateWork(1, "Pride and Prejudice", 999)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "author not found")
	})
}

func TestWorkService_DeleteWork(t *testing.T) {
	t.Parallel()
	db, cleanup := testutils.SetupTestDB(t)
	defer cleanup()

	workRepo := gorm.NewWorkRepository(db)
	writerRepo := gorm.NewWriterRepository(db)
	svc := service.NewWorkService(workRepo, writerRepo)

	work := domain.NewWork(1, "Pride and Prejudice", 1)
	require.NoError(t, workRepo.Create(work))

	err := svc.DeleteWork(1)
	require.NoError(t, err)

	_, err = workRepo.GetByID(1)
	require.Error(t, err)
}
