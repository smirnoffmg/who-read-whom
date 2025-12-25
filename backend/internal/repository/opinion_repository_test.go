package repository_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/what-writers-like/backend/internal/domain"
	"github.com/what-writers-like/backend/internal/infrastructure/database"
	"github.com/what-writers-like/backend/internal/repository"
	"github.com/what-writers-like/backend/internal/repository/gorm"
	"github.com/what-writers-like/backend/internal/testutils"
)

type testRepos struct {
	writerRepo  repository.WriterRepository
	workRepo    repository.WorkRepository
	opinionRepo repository.OpinionRepository
}

func setupTestData(
	t *testing.T,
	db *database.Database,
) (*testRepos, *domain.Writer, *domain.Writer, *domain.Work, *domain.Opinion) {
	writerRepo := gorm.NewWriterRepository(db)
	workRepo := gorm.NewWorkRepository(db)
	opinionRepo := gorm.NewOpinionRepository(db)

	writer1 := domain.NewWriter(1, "Jane Austen", 1775, nil, nil)
	writer2 := domain.NewWriter(2, "Charlotte Bronte", 1816, nil, nil)
	require.NoError(t, writerRepo.Create(writer1))
	require.NoError(t, writerRepo.Create(writer2))

	work := domain.NewWork(1, "Pride and Prejudice", 1)
	require.NoError(t, workRepo.Create(work))

	opinion := domain.NewOpinion(2, 1, true, "A delightful novel", "Personal Letters", nil, nil)
	require.NoError(t, opinionRepo.Create(opinion))

	return &testRepos{
		writerRepo:  writerRepo,
		workRepo:    workRepo,
		opinionRepo: opinionRepo,
	}, writer1, writer2, work, opinion
}

func TestOpinionRepository_Create(t *testing.T) {
	t.Parallel()
	db, cleanup := testutils.SetupTestDB(t)
	defer cleanup()

	writerRepo := gorm.NewWriterRepository(db)
	workRepo := gorm.NewWorkRepository(db)
	opinionRepo := gorm.NewOpinionRepository(db)

	writer1 := domain.NewWriter(1, "Jane Austen", 1775, nil, nil)
	writer2 := domain.NewWriter(2, "Charlotte Bronte", 1816, nil, nil)
	require.NoError(t, writerRepo.Create(writer1))
	require.NoError(t, writerRepo.Create(writer2))

	work := domain.NewWork(1, "Pride and Prejudice", 1)
	require.NoError(t, workRepo.Create(work))

	opinion := domain.NewOpinion(2, 1, true, "A delightful novel", "Personal Letters", nil, nil)
	err := opinionRepo.Create(opinion)
	assert.NoError(t, err)
}

func TestOpinionRepository_GetByWriterID(t *testing.T) {
	t.Parallel()
	db, cleanup := testutils.SetupTestDB(t)
	defer cleanup()

	repos, _, _, _, opinion := setupTestData(t, db)

	opinions, err := repos.opinionRepo.GetByWriterID(2)
	require.NoError(t, err)
	assert.Len(t, opinions, 1)
	assert.Equal(t, opinion.WriterID(), opinions[0].WriterID())
}

func TestOpinionRepository_GetByWorkID(t *testing.T) {
	t.Parallel()
	db, cleanup := testutils.SetupTestDB(t)
	defer cleanup()

	repos, _, _, _, opinion := setupTestData(t, db)

	opinions, err := repos.opinionRepo.GetByWorkID(1)
	require.NoError(t, err)
	assert.Len(t, opinions, 1)
	assert.Equal(t, opinion.WorkID(), opinions[0].WorkID())
}

func TestOpinionRepository_GetByWriterAndWork(t *testing.T) {
	t.Parallel()
	db, cleanup := testutils.SetupTestDB(t)
	defer cleanup()

	writerRepo := gorm.NewWriterRepository(db)
	workRepo := gorm.NewWorkRepository(db)
	opinionRepo := gorm.NewOpinionRepository(db)

	writer1 := domain.NewWriter(1, "Jane Austen", 1775, nil, nil)
	writer2 := domain.NewWriter(2, "Charlotte Bronte", 1816, nil, nil)
	require.NoError(t, writerRepo.Create(writer1))
	require.NoError(t, writerRepo.Create(writer2))

	work := domain.NewWork(1, "Pride and Prejudice", 1)
	require.NoError(t, workRepo.Create(work))

	opinion := domain.NewOpinion(2, 1, true, "A delightful novel", "Personal Letters", nil, nil)
	require.NoError(t, opinionRepo.Create(opinion))

	found, err := opinionRepo.GetByWriterAndWork(2, 1)
	require.NoError(t, err)
	assert.Equal(t, opinion.WriterID(), found.WriterID())
	assert.Equal(t, opinion.WorkID(), found.WorkID())
}

func TestOpinionRepository_List(t *testing.T) {
	t.Parallel()
	db, cleanup := testutils.SetupTestDB(t)
	defer cleanup()

	writerRepo := gorm.NewWriterRepository(db)
	workRepo := gorm.NewWorkRepository(db)
	opinionRepo := gorm.NewOpinionRepository(db)

	writer1 := domain.NewWriter(1, "Jane Austen", 1775, nil, nil)
	writer2 := domain.NewWriter(2, "Charlotte Bronte", 1816, nil, nil)
	require.NoError(t, writerRepo.Create(writer1))
	require.NoError(t, writerRepo.Create(writer2))

	work := domain.NewWork(1, "Pride and Prejudice", 1)
	require.NoError(t, workRepo.Create(work))

	work2 := domain.NewWork(2, "Emma", 1)
	require.NoError(t, workRepo.Create(work2))

	opinion1 := domain.NewOpinion(2, 1, true, "A delightful novel", "Personal Letters", nil, nil)
	opinion2 := domain.NewOpinion(2, 2, false, "Overrated", "Another Source", nil, nil)
	require.NoError(t, opinionRepo.Create(opinion1))
	require.NoError(t, opinionRepo.Create(opinion2))

	opinions, err := opinionRepo.List(10, 0)
	require.NoError(t, err)
	assert.Len(t, opinions, 2)
}

func TestOpinionRepository_Update(t *testing.T) {
	t.Parallel()
	db, cleanup := testutils.SetupTestDB(t)
	defer cleanup()

	writerRepo := gorm.NewWriterRepository(db)
	workRepo := gorm.NewWorkRepository(db)
	opinionRepo := gorm.NewOpinionRepository(db)

	writer1 := domain.NewWriter(1, "Jane Austen", 1775, nil, nil)
	writer2 := domain.NewWriter(2, "Charlotte Bronte", 1816, nil, nil)
	require.NoError(t, writerRepo.Create(writer1))
	require.NoError(t, writerRepo.Create(writer2))

	work := domain.NewWork(1, "Pride and Prejudice", 1)
	require.NoError(t, workRepo.Create(work))

	opinion := domain.NewOpinion(2, 1, true, "A delightful novel", "Personal Letters", nil, nil)
	require.NoError(t, opinionRepo.Create(opinion))

	updated := domain.NewOpinion(2, 1, false, "Actually, it's overrated", "Personal Letters", nil, nil)
	err := opinionRepo.Update(updated)
	require.NoError(t, err)

	found, err := opinionRepo.GetByWriterAndWork(2, 1)
	require.NoError(t, err)
	assert.False(t, found.Sentiment())
	assert.Equal(t, "Actually, it's overrated", found.Quote())
}

func TestOpinionRepository_Delete(t *testing.T) {
	t.Parallel()
	db, cleanup := testutils.SetupTestDB(t)
	defer cleanup()

	writerRepo := gorm.NewWriterRepository(db)
	workRepo := gorm.NewWorkRepository(db)
	opinionRepo := gorm.NewOpinionRepository(db)

	writer1 := domain.NewWriter(1, "Jane Austen", 1775, nil, nil)
	writer2 := domain.NewWriter(2, "Charlotte Bronte", 1816, nil, nil)
	require.NoError(t, writerRepo.Create(writer1))
	require.NoError(t, writerRepo.Create(writer2))

	work := domain.NewWork(1, "Pride and Prejudice", 1)
	require.NoError(t, workRepo.Create(work))

	opinion := domain.NewOpinion(2, 1, true, "A delightful novel", "Personal Letters", nil, nil)
	require.NoError(t, opinionRepo.Create(opinion))

	err := opinionRepo.Delete(2, 1)
	require.NoError(t, err)

	_, err = opinionRepo.GetByWriterAndWork(2, 1)
	require.Error(t, err)
}
