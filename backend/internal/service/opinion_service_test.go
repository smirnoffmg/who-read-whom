package service_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/what-writers-like/backend/internal/domain"
	"github.com/what-writers-like/backend/internal/service"
)

type mockOpinionRepository struct {
	opinions           map[string]*domain.Opinion
	create             func(*domain.Opinion) error
	getByWriterID      func(uint64) ([]*domain.Opinion, error)
	getByWorkID        func(uint64) ([]*domain.Opinion, error)
	getByWriterAndWork func(uint64, uint64) (*domain.Opinion, error)
	list               func(int, int) ([]*domain.Opinion, error)
	update             func(*domain.Opinion) error
	delete             func(uint64, uint64) error
}

func (m *mockOpinionRepository) Create(opinion *domain.Opinion) error {
	if m.create != nil {
		return m.create(opinion)
	}
	if m.opinions == nil {
		m.opinions = make(map[string]*domain.Opinion)
	}
	key := key(opinion.WriterID(), opinion.WorkID())
	m.opinions[key] = opinion
	return nil
}

func (m *mockOpinionRepository) GetByWriterID(writerID uint64) ([]*domain.Opinion, error) {
	if m.getByWriterID != nil {
		return m.getByWriterID(writerID)
	}
	if m.opinions == nil {
		return []*domain.Opinion{}, nil
	}
	result := make([]*domain.Opinion, 0)
	for _, o := range m.opinions {
		if o.WriterID() == writerID {
			result = append(result, o)
		}
	}
	return result, nil
}

func (m *mockOpinionRepository) GetByWorkID(workID uint64) ([]*domain.Opinion, error) {
	if m.getByWorkID != nil {
		return m.getByWorkID(workID)
	}
	if m.opinions == nil {
		return []*domain.Opinion{}, nil
	}
	result := make([]*domain.Opinion, 0)
	for _, o := range m.opinions {
		if o.WorkID() == workID {
			result = append(result, o)
		}
	}
	return result, nil
}

func (m *mockOpinionRepository) GetByWriterAndWork(writerID, workID uint64) (*domain.Opinion, error) {
	if m.getByWriterAndWork != nil {
		return m.getByWriterAndWork(writerID, workID)
	}
	if m.opinions == nil {
		return nil, errors.New("not found")
	}
	key := key(writerID, workID)
	opinion, ok := m.opinions[key]
	if !ok {
		return nil, errors.New("not found")
	}
	return opinion, nil
}

func (m *mockOpinionRepository) List(limit, offset int) ([]*domain.Opinion, error) {
	if m.list != nil {
		return m.list(limit, offset)
	}
	if m.opinions == nil {
		return []*domain.Opinion{}, nil
	}
	result := make([]*domain.Opinion, 0, len(m.opinions))
	for _, o := range m.opinions {
		result = append(result, o)
	}
	return result, nil
}

func (m *mockOpinionRepository) Update(opinion *domain.Opinion) error {
	if m.update != nil {
		return m.update(opinion)
	}
	if m.opinions == nil {
		m.opinions = make(map[string]*domain.Opinion)
	}
	key := key(opinion.WriterID(), opinion.WorkID())
	m.opinions[key] = opinion
	return nil
}

func (m *mockOpinionRepository) Delete(writerID, workID uint64) error {
	if m.delete != nil {
		return m.delete(writerID, workID)
	}
	if m.opinions == nil {
		return errors.New("not found")
	}
	key := key(writerID, workID)
	delete(m.opinions, key)
	return nil
}

func key(writerID, workID uint64) string {
	return string(rune(writerID)) + ":" + string(rune(workID))
}

func TestOpinionService_ListOpinions(t *testing.T) {
	t.Parallel()
	expectedOpinions := []*domain.Opinion{
		domain.NewOpinion(2, 1, true, "Quote 1", "Source 1", nil, nil),
		domain.NewOpinion(3, 2, false, "Quote 2", "Source 2", nil, nil),
	}
	opinionRepo := &mockOpinionRepository{
		list: func(int, int) ([]*domain.Opinion, error) {
			return expectedOpinions, nil
		},
	}
	writerRepo := &mockWriterRepository{}
	workRepo := &mockWorkRepository{}
	svc := service.NewOpinionService(opinionRepo, writerRepo, workRepo)

	opinions, err := svc.ListOpinions(10, 0)
	require.NoError(t, err)
	assert.Len(t, opinions, 2)
}

func TestOpinionService_CreateOpinion(t *testing.T) {
	t.Parallel()
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		opinionRepo := &mockOpinionRepository{}
		writerRepo := &mockWriterRepository{
			getByID: func(id uint64) (*domain.Writer, error) {
				if id == 2 {
					return domain.NewWriter(2, "Charlotte Bronte", 1816, nil, nil), nil
				}
				return nil, errors.New("not found")
			},
		}
		workRepo := &mockWorkRepository{
			works: map[uint64]*domain.Work{
				1: domain.NewWork(1, "Pride and Prejudice", 1),
			},
		}
		svc := service.NewOpinionService(opinionRepo, writerRepo, workRepo)

		opinion, err := svc.CreateOpinion(2, 1, true, "A delightful novel", "Personal Letters", nil, nil)
		require.NoError(t, err)
		assert.Equal(t, uint64(2), opinion.WriterID())
		assert.Equal(t, uint64(1), opinion.WorkID())
		assert.True(t, opinion.Sentiment())
	})

	t.Run("empty quote", func(t *testing.T) {
		t.Parallel()
		opinionRepo := &mockOpinionRepository{}
		writerRepo := &mockWriterRepository{}
		workRepo := &mockWorkRepository{}
		svc := service.NewOpinionService(opinionRepo, writerRepo, workRepo)

		_, err := svc.CreateOpinion(2, 1, true, "", "Source", nil, nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "quote is required")
	})

	t.Run("empty source", func(t *testing.T) {
		t.Parallel()
		opinionRepo := &mockOpinionRepository{}
		writerRepo := &mockWriterRepository{}
		workRepo := &mockWorkRepository{}
		svc := service.NewOpinionService(opinionRepo, writerRepo, workRepo)

		_, err := svc.CreateOpinion(2, 1, true, "Quote", "", nil, nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "source is required")
	})

	t.Run("work not found", func(t *testing.T) {
		t.Parallel()
		opinionRepo := &mockOpinionRepository{}
		writerRepo := &mockWriterRepository{}
		workRepo := &mockWorkRepository{}
		svc := service.NewOpinionService(opinionRepo, writerRepo, workRepo)

		_, err := svc.CreateOpinion(2, 999, true, "Quote", "Source", nil, nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "work not found")
	})

	t.Run("writer cannot express opinion about own work", func(t *testing.T) {
		t.Parallel()
		opinionRepo := &mockOpinionRepository{}
		writerRepo := &mockWriterRepository{
			getByID: func(id uint64) (*domain.Writer, error) {
				if id == 1 {
					return domain.NewWriter(1, "Jane Austen", 1775, nil, nil), nil
				}
				return nil, errors.New("not found")
			},
		}
		workRepo := &mockWorkRepository{
			works: map[uint64]*domain.Work{
				1: domain.NewWork(1, "Pride and Prejudice", 1),
			},
		}
		svc := service.NewOpinionService(opinionRepo, writerRepo, workRepo)

		_, err := svc.CreateOpinion(1, 1, true, "Quote", "Source", nil, nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "writer cannot express opinion about their own work")
	})

	t.Run("writer not found", func(t *testing.T) {
		t.Parallel()
		opinionRepo := &mockOpinionRepository{}
		writerRepo := &mockWriterRepository{
			getByID: func(uint64) (*domain.Writer, error) {
				return nil, errors.New("not found")
			},
		}
		workRepo := &mockWorkRepository{
			works: map[uint64]*domain.Work{
				1: domain.NewWork(1, "Pride and Prejudice", 1),
			},
		}
		svc := service.NewOpinionService(opinionRepo, writerRepo, workRepo)

		_, err := svc.CreateOpinion(999, 1, true, "Quote", "Source", nil, nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "writer not found")
	})
}

func TestOpinionService_GetOpinionsByWriter(t *testing.T) {
	t.Parallel()
	expectedOpinions := []*domain.Opinion{
		domain.NewOpinion(2, 1, true, "Quote 1", "Source 1", nil, nil),
		domain.NewOpinion(2, 2, false, "Quote 2", "Source 2", nil, nil),
	}
	opinionRepo := &mockOpinionRepository{
		getByWriterID: func(uint64) ([]*domain.Opinion, error) {
			return expectedOpinions, nil
		},
	}
	writerRepo := &mockWriterRepository{}
	workRepo := &mockWorkRepository{}
	svc := service.NewOpinionService(opinionRepo, writerRepo, workRepo)

	opinions, err := svc.GetOpinionsByWriter(2)
	require.NoError(t, err)
	assert.Len(t, opinions, 2)
}

func TestOpinionService_GetOpinionsByWork(t *testing.T) {
	t.Parallel()
	expectedOpinions := []*domain.Opinion{
		domain.NewOpinion(2, 1, true, "Quote 1", "Source 1", nil, nil),
		domain.NewOpinion(3, 1, false, "Quote 2", "Source 2", nil, nil),
	}
	opinionRepo := &mockOpinionRepository{
		getByWorkID: func(uint64) ([]*domain.Opinion, error) {
			return expectedOpinions, nil
		},
	}
	writerRepo := &mockWriterRepository{}
	workRepo := &mockWorkRepository{}
	svc := service.NewOpinionService(opinionRepo, writerRepo, workRepo)

	opinions, err := svc.GetOpinionsByWork(1)
	require.NoError(t, err)
	assert.Len(t, opinions, 2)
}

func TestOpinionService_GetOpinion(t *testing.T) {
	t.Parallel()
	expectedOpinion := domain.NewOpinion(2, 1, true, "Quote", "Source", nil, nil)
	opinionRepo := &mockOpinionRepository{
		getByWriterAndWork: func(writerID, workID uint64) (*domain.Opinion, error) {
			if writerID == 2 && workID == 1 {
				return expectedOpinion, nil
			}
			return nil, errors.New("not found")
		},
	}
	writerRepo := &mockWriterRepository{}
	workRepo := &mockWorkRepository{}
	svc := service.NewOpinionService(opinionRepo, writerRepo, workRepo)

	opinion, err := svc.GetOpinion(2, 1)
	require.NoError(t, err)
	assert.Equal(t, expectedOpinion.WriterID(), opinion.WriterID())
	assert.Equal(t, expectedOpinion.WorkID(), opinion.WorkID())
}

func TestOpinionService_UpdateOpinion(t *testing.T) {
	t.Parallel()
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		opinionRepo := &mockOpinionRepository{
			update: func(*domain.Opinion) error {
				return nil
			},
		}
		writerRepo := &mockWriterRepository{}
		workRepo := &mockWorkRepository{
			works: map[uint64]*domain.Work{
				1: domain.NewWork(1, "Pride and Prejudice", 1),
			},
		}
		svc := service.NewOpinionService(opinionRepo, writerRepo, workRepo)

		err := svc.UpdateOpinion(2, 1, false, "Updated quote", "Updated source", nil, nil)
		assert.NoError(t, err)
	})

	t.Run("empty quote", func(t *testing.T) {
		t.Parallel()
		opinionRepo := &mockOpinionRepository{}
		writerRepo := &mockWriterRepository{}
		workRepo := &mockWorkRepository{}
		svc := service.NewOpinionService(opinionRepo, writerRepo, workRepo)

		err := svc.UpdateOpinion(2, 1, true, "", "Source", nil, nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "quote is required")
	})

	t.Run("writer cannot express opinion about own work", func(t *testing.T) {
		t.Parallel()
		opinionRepo := &mockOpinionRepository{}
		writerRepo := &mockWriterRepository{}
		workRepo := &mockWorkRepository{
			works: map[uint64]*domain.Work{
				1: domain.NewWork(1, "Pride and Prejudice", 1),
			},
		}
		svc := service.NewOpinionService(opinionRepo, writerRepo, workRepo)

		err := svc.UpdateOpinion(1, 1, true, "Quote", "Source", nil, nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "writer cannot express opinion about their own work")
	})
}

func TestOpinionService_DeleteOpinion(t *testing.T) {
	t.Parallel()
	opinionRepo := &mockOpinionRepository{
		delete: func(uint64, uint64) error {
			return nil
		},
	}
	writerRepo := &mockWriterRepository{}
	workRepo := &mockWorkRepository{}
	svc := service.NewOpinionService(opinionRepo, writerRepo, workRepo)

	err := svc.DeleteOpinion(2, 1)
	assert.NoError(t, err)
}
