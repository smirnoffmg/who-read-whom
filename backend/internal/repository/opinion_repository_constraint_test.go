package repository_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/what-writers-like/backend/internal/domain"
	"github.com/what-writers-like/backend/internal/repository/gorm"
	"github.com/what-writers-like/backend/internal/testutils"
)

func TestOpinionRepository_DatabaseConstraint(t *testing.T) {
	t.Parallel()
	db, cleanup := testutils.SetupTestDB(t)
	defer cleanup()

	writerRepo := gorm.NewWriterRepository(db)
	workRepo := gorm.NewWorkRepository(db)
	opinionRepo := gorm.NewOpinionRepository(db)

	// Create a writer
	writer := domain.NewWriter(1, "Jane Austen", 1775, nil, nil)
	err := writerRepo.Create(writer)
	require.NoError(t, err)

	// Create a work by that writer
	work := domain.NewWork(1, "Pride and Prejudice", 1)
	err = workRepo.Create(work)
	require.NoError(t, err)

	// Try to create an opinion where writer_id = work.author_id (should fail at DB level)
	opinion := domain.NewOpinion(1, 1, true, "My own work", "Personal", nil, nil)
	err = opinionRepo.Create(opinion)

	// Should fail due to database constraint
	require.Error(t, err)
	assert.Contains(t, err.Error(), "writer cannot express opinion about their own work")
}
