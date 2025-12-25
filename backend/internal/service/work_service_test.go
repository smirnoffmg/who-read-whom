package service_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/what-writers-like/backend/internal/domain"
	"github.com/what-writers-like/backend/internal/service"
)

func TestWorkService_CreateWork(t *testing.T) {
	t.Parallel()
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		workRepo := &mockWorkRepository{}
		writerRepo := &mockWriterRepository{
			getByID: func(id uint64) (*domain.Writer, error) {
				if id == 1 {
					return domain.NewWriter(1, "Jane Austen", 1775, nil, nil), nil
				}
				return nil, errors.New("not found")
			},
		}
		svc := service.NewWorkService(workRepo, writerRepo)

		work, err := svc.CreateWork("Pride and Prejudice", 1)
		require.NoError(t, err)
		assert.Equal(t, "Pride and Prejudice", work.Title())
		assert.Equal(t, uint64(1), work.AuthorID())
	})

	t.Run("empty title", func(t *testing.T) {
		t.Parallel()
		workRepo := &mockWorkRepository{}
		writerRepo := &mockWriterRepository{}
		svc := service.NewWorkService(workRepo, writerRepo)

		_, err := svc.CreateWork("", 1)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "title is required")
	})

	t.Run("author not found", func(t *testing.T) {
		t.Parallel()
		workRepo := &mockWorkRepository{}
		writerRepo := &mockWriterRepository{
			getByID: func(uint64) (*domain.Writer, error) {
				return nil, errors.New("not found")
			},
		}
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
		expectedWork := domain.NewWork(1, "Pride and Prejudice", 1)
		workRepo := &mockWorkRepository{
			works: map[uint64]*domain.Work{
				1: expectedWork,
			},
		}
		writerRepo := &mockWriterRepository{}
		svc := service.NewWorkService(workRepo, writerRepo)

		work, err := svc.GetWork(1)
		require.NoError(t, err)
		assert.Equal(t, expectedWork.ID(), work.ID())
		assert.Equal(t, expectedWork.Title(), work.Title())
	})

	t.Run("not found", func(t *testing.T) {
		t.Parallel()
		workRepo := &mockWorkRepository{}
		writerRepo := &mockWriterRepository{}
		svc := service.NewWorkService(workRepo, writerRepo)

		_, err := svc.GetWork(999)
		require.Error(t, err)
	})
}

func TestWorkService_GetWorksByAuthor(t *testing.T) {
	t.Parallel()
	expectedWorks := []*domain.Work{
		domain.NewWork(1, "Pride and Prejudice", 1),
		domain.NewWork(2, "Sense and Sensibility", 1),
	}
	workRepo := &mockWorkRepository{
		getByAuthorID: func(uint64) ([]*domain.Work, error) {
			return expectedWorks, nil
		},
	}
	writerRepo := &mockWriterRepository{}
	svc := service.NewWorkService(workRepo, writerRepo)

	works, err := svc.GetWorksByAuthor(1)
	require.NoError(t, err)
	assert.Len(t, works, 2)
}

func TestWorkService_ListWorks(t *testing.T) {
	t.Parallel()
	expectedWorks := []*domain.Work{
		domain.NewWork(1, "Pride and Prejudice", 1),
		domain.NewWork(2, "Sense and Sensibility", 1),
	}
	workRepo := &mockWorkRepository{
		works: map[uint64]*domain.Work{
			1: expectedWorks[0],
			2: expectedWorks[1],
		},
	}
	writerRepo := &mockWriterRepository{}
	svc := service.NewWorkService(workRepo, writerRepo)

	works, err := svc.ListWorks(10, 0)
	require.NoError(t, err)
	assert.Len(t, works, 2)
}

func TestWorkService_UpdateWork(t *testing.T) {
	t.Parallel()
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		workRepo := &mockWorkRepository{
			works: make(map[uint64]*domain.Work),
		}
		writerRepo := &mockWriterRepository{
			getByID: func(id uint64) (*domain.Writer, error) {
				if id == 1 {
					return domain.NewWriter(1, "Jane Austen", 1775, nil, nil), nil
				}
				return nil, errors.New("not found")
			},
		}
		svc := service.NewWorkService(workRepo, writerRepo)

		err := svc.UpdateWork(1, "Pride and Prejudice (Revised)", 1)
		assert.NoError(t, err)
	})

	t.Run("empty title", func(t *testing.T) {
		t.Parallel()
		workRepo := &mockWorkRepository{}
		writerRepo := &mockWriterRepository{}
		svc := service.NewWorkService(workRepo, writerRepo)

		err := svc.UpdateWork(1, "", 1)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "title is required")
	})

	t.Run("author not found", func(t *testing.T) {
		t.Parallel()
		workRepo := &mockWorkRepository{}
		writerRepo := &mockWriterRepository{
			getByID: func(uint64) (*domain.Writer, error) {
				return nil, errors.New("not found")
			},
		}
		svc := service.NewWorkService(workRepo, writerRepo)

		err := svc.UpdateWork(1, "Pride and Prejudice", 999)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "author not found")
	})
}

func TestWorkService_DeleteWork(t *testing.T) {
	t.Parallel()
	workRepo := &mockWorkRepository{
		works: map[uint64]*domain.Work{
			1: domain.NewWork(1, "Pride and Prejudice", 1),
		},
		delete: func(uint64) error {
			return nil
		},
	}
	writerRepo := &mockWriterRepository{}
	svc := service.NewWorkService(workRepo, writerRepo)

	err := svc.DeleteWork(1)
	assert.NoError(t, err)
}
