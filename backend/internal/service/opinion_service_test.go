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

func TestOpinionService_ListOpinions(t *testing.T) {
	t.Parallel()
	db, cleanup := testutils.SetupTestDB(t)
	defer cleanup()

	opinionRepo := gorm.NewOpinionRepository(db)
	writerRepo := gorm.NewWriterRepository(db)
	workRepo := gorm.NewWorkRepository(db)
	svc := service.NewOpinionService(opinionRepo, writerRepo, workRepo)

	writer1 := domain.NewWriter(1, "Jane Austen", 1775, nil, nil)
	writer2 := domain.NewWriter(2, "Charlotte Bronte", 1816, nil, nil)
	writer3 := domain.NewWriter(3, "Charles Dickens", 1812, nil, nil)
	require.NoError(t, writerRepo.Create(writer1))
	require.NoError(t, writerRepo.Create(writer2))
	require.NoError(t, writerRepo.Create(writer3))

	work1 := domain.NewWork(1, "Pride and Prejudice", 1)
	work2 := domain.NewWork(2, "Jane Eyre", 2)
	require.NoError(t, workRepo.Create(work1))
	require.NoError(t, workRepo.Create(work2))

	opinion1 := domain.NewOpinion(2, 1, true, "Quote 1", "Source 1", nil, nil)
	opinion2 := domain.NewOpinion(3, 2, false, "Quote 2", "Source 2", nil, nil)
	require.NoError(t, opinionRepo.Create(opinion1))
	require.NoError(t, opinionRepo.Create(opinion2))

	opinions, err := svc.ListOpinions(10, 0)
	require.NoError(t, err)
	assert.Len(t, opinions, 2)
}

func TestOpinionService_CreateOpinion(t *testing.T) {
	t.Parallel()
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		db, cleanup := testutils.SetupTestDB(t)
		defer cleanup()

		opinionRepo := gorm.NewOpinionRepository(db)
		writerRepo := gorm.NewWriterRepository(db)
		workRepo := gorm.NewWorkRepository(db)
		svc := service.NewOpinionService(opinionRepo, writerRepo, workRepo)

		writer1 := domain.NewWriter(1, "Jane Austen", 1775, nil, nil)
		writer2 := domain.NewWriter(2, "Charlotte Bronte", 1816, nil, nil)
		require.NoError(t, writerRepo.Create(writer1))
		require.NoError(t, writerRepo.Create(writer2))

		work := domain.NewWork(1, "Pride and Prejudice", 1)
		require.NoError(t, workRepo.Create(work))

		opinion, err := svc.CreateOpinion(2, 1, true, "A delightful novel", "Personal Letters", nil, nil)
		require.NoError(t, err)
		assert.Equal(t, uint64(2), opinion.WriterID())
		assert.Equal(t, uint64(1), opinion.WorkID())
		assert.True(t, opinion.Sentiment())
	})

	t.Run("empty quote", func(t *testing.T) {
		t.Parallel()
		db, cleanup := testutils.SetupTestDB(t)
		defer cleanup()

		opinionRepo := gorm.NewOpinionRepository(db)
		writerRepo := gorm.NewWriterRepository(db)
		workRepo := gorm.NewWorkRepository(db)
		svc := service.NewOpinionService(opinionRepo, writerRepo, workRepo)

		_, err := svc.CreateOpinion(2, 1, true, "", "Source", nil, nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "quote is required")
	})

	t.Run("empty source", func(t *testing.T) {
		t.Parallel()
		db, cleanup := testutils.SetupTestDB(t)
		defer cleanup()

		opinionRepo := gorm.NewOpinionRepository(db)
		writerRepo := gorm.NewWriterRepository(db)
		workRepo := gorm.NewWorkRepository(db)
		svc := service.NewOpinionService(opinionRepo, writerRepo, workRepo)

		_, err := svc.CreateOpinion(2, 1, true, "Quote", "", nil, nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "source is required")
	})

	t.Run("work not found", func(t *testing.T) {
		t.Parallel()
		db, cleanup := testutils.SetupTestDB(t)
		defer cleanup()

		opinionRepo := gorm.NewOpinionRepository(db)
		writerRepo := gorm.NewWriterRepository(db)
		workRepo := gorm.NewWorkRepository(db)
		svc := service.NewOpinionService(opinionRepo, writerRepo, workRepo)

		_, err := svc.CreateOpinion(2, 999, true, "Quote", "Source", nil, nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "work not found")
	})

	t.Run("writer cannot express opinion about own work", func(t *testing.T) {
		t.Parallel()
		db, cleanup := testutils.SetupTestDB(t)
		defer cleanup()

		opinionRepo := gorm.NewOpinionRepository(db)
		writerRepo := gorm.NewWriterRepository(db)
		workRepo := gorm.NewWorkRepository(db)
		svc := service.NewOpinionService(opinionRepo, writerRepo, workRepo)

		writer := domain.NewWriter(1, "Jane Austen", 1775, nil, nil)
		require.NoError(t, writerRepo.Create(writer))
		work := domain.NewWork(1, "Pride and Prejudice", 1)
		require.NoError(t, workRepo.Create(work))

		_, err := svc.CreateOpinion(1, 1, true, "Quote", "Source", nil, nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "writer cannot express opinion about their own work")
	})

	t.Run("writer not found", func(t *testing.T) {
		t.Parallel()
		db, cleanup := testutils.SetupTestDB(t)
		defer cleanup()

		opinionRepo := gorm.NewOpinionRepository(db)
		writerRepo := gorm.NewWriterRepository(db)
		workRepo := gorm.NewWorkRepository(db)
		svc := service.NewOpinionService(opinionRepo, writerRepo, workRepo)

		work := domain.NewWork(1, "Pride and Prejudice", 1)
		require.NoError(t, workRepo.Create(work))

		_, err := svc.CreateOpinion(999, 1, true, "Quote", "Source", nil, nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "writer not found")
	})
}

func TestOpinionService_GetOpinionsByWriter(t *testing.T) {
	t.Parallel()
	db, cleanup := testutils.SetupTestDB(t)
	defer cleanup()

	opinionRepo := gorm.NewOpinionRepository(db)
	writerRepo := gorm.NewWriterRepository(db)
	workRepo := gorm.NewWorkRepository(db)
	svc := service.NewOpinionService(opinionRepo, writerRepo, workRepo)

	writer1 := domain.NewWriter(1, "Jane Austen", 1775, nil, nil)
	writer2 := domain.NewWriter(2, "Charlotte Bronte", 1816, nil, nil)
	writer3 := domain.NewWriter(3, "Charles Dickens", 1812, nil, nil)
	require.NoError(t, writerRepo.Create(writer1))
	require.NoError(t, writerRepo.Create(writer2))
	require.NoError(t, writerRepo.Create(writer3))

	work1 := domain.NewWork(1, "Pride and Prejudice", 1)
	work2 := domain.NewWork(2, "Jane Eyre", 2)
	require.NoError(t, workRepo.Create(work1))
	require.NoError(t, workRepo.Create(work2))

	// Writer 2 (Charlotte Bronte) expresses opinion about work 1 (Jane Austen's work)
	// Writer 3 (Charles Dickens) expresses opinion about work 2 (Charlotte Bronte's work)
	opinion1 := domain.NewOpinion(2, 1, true, "Quote 1", "Source 1", nil, nil)
	opinion2 := domain.NewOpinion(3, 2, false, "Quote 2", "Source 2", nil, nil)
	require.NoError(t, opinionRepo.Create(opinion1))
	require.NoError(t, opinionRepo.Create(opinion2))

	opinions, err := svc.GetOpinionsByWriter(2)
	require.NoError(t, err)
	assert.Len(t, opinions, 1)
	assert.Equal(t, uint64(1), opinions[0].WorkID())
}

func TestOpinionService_GetOpinionsByWork(t *testing.T) {
	t.Parallel()
	db, cleanup := testutils.SetupTestDB(t)
	defer cleanup()

	opinionRepo := gorm.NewOpinionRepository(db)
	writerRepo := gorm.NewWriterRepository(db)
	workRepo := gorm.NewWorkRepository(db)
	svc := service.NewOpinionService(opinionRepo, writerRepo, workRepo)

	writer1 := domain.NewWriter(1, "Jane Austen", 1775, nil, nil)
	writer2 := domain.NewWriter(2, "Charlotte Bronte", 1816, nil, nil)
	writer3 := domain.NewWriter(3, "Charles Dickens", 1812, nil, nil)
	require.NoError(t, writerRepo.Create(writer1))
	require.NoError(t, writerRepo.Create(writer2))
	require.NoError(t, writerRepo.Create(writer3))

	work := domain.NewWork(1, "Pride and Prejudice", 1)
	require.NoError(t, workRepo.Create(work))

	opinion1 := domain.NewOpinion(2, 1, true, "Quote 1", "Source 1", nil, nil)
	opinion2 := domain.NewOpinion(3, 1, false, "Quote 2", "Source 2", nil, nil)
	require.NoError(t, opinionRepo.Create(opinion1))
	require.NoError(t, opinionRepo.Create(opinion2))

	opinions, err := svc.GetOpinionsByWork(1)
	require.NoError(t, err)
	assert.Len(t, opinions, 2)
}

func TestOpinionService_GetOpinion(t *testing.T) {
	t.Parallel()
	db, cleanup := testutils.SetupTestDB(t)
	defer cleanup()

	opinionRepo := gorm.NewOpinionRepository(db)
	writerRepo := gorm.NewWriterRepository(db)
	workRepo := gorm.NewWorkRepository(db)
	svc := service.NewOpinionService(opinionRepo, writerRepo, workRepo)

	writer1 := domain.NewWriter(1, "Jane Austen", 1775, nil, nil)
	writer2 := domain.NewWriter(2, "Charlotte Bronte", 1816, nil, nil)
	require.NoError(t, writerRepo.Create(writer1))
	require.NoError(t, writerRepo.Create(writer2))

	work := domain.NewWork(1, "Pride and Prejudice", 1)
	require.NoError(t, workRepo.Create(work))

	expectedOpinion := domain.NewOpinion(2, 1, true, "Quote", "Source", nil, nil)
	require.NoError(t, opinionRepo.Create(expectedOpinion))

	opinion, err := svc.GetOpinion(2, 1)
	require.NoError(t, err)
	assert.Equal(t, expectedOpinion.WriterID(), opinion.WriterID())
	assert.Equal(t, expectedOpinion.WorkID(), opinion.WorkID())
}

func TestOpinionService_UpdateOpinion(t *testing.T) {
	t.Parallel()
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		db, cleanup := testutils.SetupTestDB(t)
		defer cleanup()

		opinionRepo := gorm.NewOpinionRepository(db)
		writerRepo := gorm.NewWriterRepository(db)
		workRepo := gorm.NewWorkRepository(db)
		svc := service.NewOpinionService(opinionRepo, writerRepo, workRepo)

		writer1 := domain.NewWriter(1, "Jane Austen", 1775, nil, nil)
		writer2 := domain.NewWriter(2, "Charlotte Bronte", 1816, nil, nil)
		require.NoError(t, writerRepo.Create(writer1))
		require.NoError(t, writerRepo.Create(writer2))

		work := domain.NewWork(1, "Pride and Prejudice", 1)
		require.NoError(t, workRepo.Create(work))

		opinion := domain.NewOpinion(2, 1, true, "Quote", "Source", nil, nil)
		require.NoError(t, opinionRepo.Create(opinion))

		err := svc.UpdateOpinion(2, 1, false, "Updated quote", "Updated source", nil, nil)
		require.NoError(t, err)

		updated, err := opinionRepo.GetByWriterAndWork(2, 1)
		require.NoError(t, err)
		assert.False(t, updated.Sentiment())
		assert.Equal(t, "Updated quote", updated.Quote())
	})

	t.Run("empty quote", func(t *testing.T) {
		t.Parallel()
		db, cleanup := testutils.SetupTestDB(t)
		defer cleanup()

		opinionRepo := gorm.NewOpinionRepository(db)
		writerRepo := gorm.NewWriterRepository(db)
		workRepo := gorm.NewWorkRepository(db)
		svc := service.NewOpinionService(opinionRepo, writerRepo, workRepo)

		err := svc.UpdateOpinion(2, 1, true, "", "Source", nil, nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "quote is required")
	})

	t.Run("writer cannot express opinion about own work", func(t *testing.T) {
		t.Parallel()
		db, cleanup := testutils.SetupTestDB(t)
		defer cleanup()

		opinionRepo := gorm.NewOpinionRepository(db)
		writerRepo := gorm.NewWriterRepository(db)
		workRepo := gorm.NewWorkRepository(db)
		svc := service.NewOpinionService(opinionRepo, writerRepo, workRepo)

		writer := domain.NewWriter(1, "Jane Austen", 1775, nil, nil)
		require.NoError(t, writerRepo.Create(writer))
		work := domain.NewWork(1, "Pride and Prejudice", 1)
		require.NoError(t, workRepo.Create(work))

		err := svc.UpdateOpinion(1, 1, true, "Quote", "Source", nil, nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "writer cannot express opinion about their own work")
	})
}

func TestOpinionService_DeleteOpinion(t *testing.T) {
	t.Parallel()
	db, cleanup := testutils.SetupTestDB(t)
	defer cleanup()

	opinionRepo := gorm.NewOpinionRepository(db)
	writerRepo := gorm.NewWriterRepository(db)
	workRepo := gorm.NewWorkRepository(db)
	svc := service.NewOpinionService(opinionRepo, writerRepo, workRepo)

	writer1 := domain.NewWriter(1, "Jane Austen", 1775, nil, nil)
	writer2 := domain.NewWriter(2, "Charlotte Bronte", 1816, nil, nil)
	require.NoError(t, writerRepo.Create(writer1))
	require.NoError(t, writerRepo.Create(writer2))

	work := domain.NewWork(1, "Pride and Prejudice", 1)
	require.NoError(t, workRepo.Create(work))

	opinion := domain.NewOpinion(2, 1, true, "Quote", "Source", nil, nil)
	require.NoError(t, opinionRepo.Create(opinion))

	err := svc.DeleteOpinion(2, 1)
	require.NoError(t, err)

	_, err = opinionRepo.GetByWriterAndWork(2, 1)
	require.Error(t, err)
}
