package service_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/what-writers-like/backend/internal/domain"
	"github.com/what-writers-like/backend/internal/service"
)

type mockWriterRepository struct {
	writers map[uint64]*domain.Writer
	create  func(*domain.Writer) error
	getByID func(uint64) (*domain.Writer, error)
	list    func(int, int) ([]*domain.Writer, error)
	update  func(*domain.Writer) error
	delete  func(uint64) error
}

func (m *mockWriterRepository) Create(writer *domain.Writer) error {
	if m.create != nil {
		return m.create(writer)
	}
	if m.writers == nil {
		m.writers = make(map[uint64]*domain.Writer)
	}
	m.writers[writer.ID()] = writer
	return nil
}

func (m *mockWriterRepository) GetByID(id uint64) (*domain.Writer, error) {
	if m.getByID != nil {
		return m.getByID(id)
	}
	if m.writers == nil {
		return nil, errors.New("not found")
	}
	writer, ok := m.writers[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return writer, nil
}

func (m *mockWriterRepository) List(limit, offset int) ([]*domain.Writer, error) {
	if m.list != nil {
		return m.list(limit, offset)
	}
	if m.writers == nil {
		return []*domain.Writer{}, nil
	}
	result := make([]*domain.Writer, 0, len(m.writers))
	for _, w := range m.writers {
		result = append(result, w)
	}
	return result, nil
}

func (m *mockWriterRepository) Update(writer *domain.Writer) error {
	if m.update != nil {
		return m.update(writer)
	}
	if m.writers == nil {
		m.writers = make(map[uint64]*domain.Writer)
	}
	m.writers[writer.ID()] = writer
	return nil
}

func (m *mockWriterRepository) Delete(id uint64) error {
	if m.delete != nil {
		return m.delete(id)
	}
	if m.writers == nil {
		return errors.New("not found")
	}
	delete(m.writers, id)
	return nil
}

type mockWorkRepository struct {
	works         map[uint64]*domain.Work
	getByAuthorID func(uint64) ([]*domain.Work, error)
	update        func(*domain.Work) error
	delete        func(uint64) error
}

func (m *mockWorkRepository) Create(work *domain.Work) error {
	if m.works == nil {
		m.works = make(map[uint64]*domain.Work)
	}
	m.works[work.ID()] = work
	return nil
}

func (m *mockWorkRepository) GetByID(id uint64) (*domain.Work, error) {
	if m.works == nil {
		return nil, errors.New("not found")
	}
	work, ok := m.works[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return work, nil
}

func (m *mockWorkRepository) GetByAuthorID(authorID uint64) ([]*domain.Work, error) {
	if m.getByAuthorID != nil {
		return m.getByAuthorID(authorID)
	}
	if m.works == nil {
		return []*domain.Work{}, nil
	}
	result := make([]*domain.Work, 0)
	for _, w := range m.works {
		if w.AuthorID() == authorID {
			result = append(result, w)
		}
	}
	return result, nil
}

func (m *mockWorkRepository) List(limit, offset int) ([]*domain.Work, error) {
	if m.works == nil {
		return []*domain.Work{}, nil
	}
	result := make([]*domain.Work, 0, len(m.works))
	for _, w := range m.works {
		result = append(result, w)
	}
	return result, nil
}

func (m *mockWorkRepository) Update(work *domain.Work) error {
	if m.update != nil {
		return m.update(work)
	}
	if m.works == nil {
		m.works = make(map[uint64]*domain.Work)
	}
	m.works[work.ID()] = work
	return nil
}

func (m *mockWorkRepository) Delete(id uint64) error {
	if m.delete != nil {
		return m.delete(id)
	}
	if m.works == nil {
		return errors.New("not found")
	}
	delete(m.works, id)
	return nil
}

func TestWriterService_CreateWriter(t *testing.T) {
	t.Parallel()
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		writerRepo := &mockWriterRepository{}
		workRepo := &mockWorkRepository{}
		svc := service.NewWriterService(writerRepo, workRepo)

		writer, err := svc.CreateWriter("Jane Austen", 1775, nil, nil)
		require.NoError(t, err)
		assert.Equal(t, "Jane Austen", writer.Name())
		assert.Equal(t, 1775, writer.BirthYear())
	})

	t.Run("empty name", func(t *testing.T) {
		t.Parallel()
		writerRepo := &mockWriterRepository{}
		workRepo := &mockWorkRepository{}
		svc := service.NewWriterService(writerRepo, workRepo)

		_, err := svc.CreateWriter("", 1775, nil, nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "name is required")
	})

	t.Run("invalid birth year", func(t *testing.T) {
		t.Parallel()
		writerRepo := &mockWriterRepository{}
		workRepo := &mockWorkRepository{}
		svc := service.NewWriterService(writerRepo, workRepo)

		_, err := svc.CreateWriter("Jane Austen", 0, nil, nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "birth year must be positive")
	})

	t.Run("increments id correctly", func(t *testing.T) {
		t.Parallel()
		writerRepo := &mockWriterRepository{
			writers: map[uint64]*domain.Writer{
				1: domain.NewWriter(1, "Writer 1", 1800, nil, nil),
				3: domain.NewWriter(3, "Writer 3", 1800, nil, nil),
			},
		}
		workRepo := &mockWorkRepository{}
		svc := service.NewWriterService(writerRepo, workRepo)

		writer, err := svc.CreateWriter("New Writer", 1900, nil, nil)
		require.NoError(t, err)
		assert.Equal(t, uint64(4), writer.ID())
	})
}

func TestWriterService_GetWriter(t *testing.T) {
	t.Parallel()
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		expectedWriter := domain.NewWriter(1, "Jane Austen", 1775, nil, nil)
		writerRepo := &mockWriterRepository{
			getByID: func(id uint64) (*domain.Writer, error) {
				if id == 1 {
					return expectedWriter, nil
				}
				return nil, errors.New("not found")
			},
		}
		workRepo := &mockWorkRepository{}
		svc := service.NewWriterService(writerRepo, workRepo)

		writer, err := svc.GetWriter(1)
		require.NoError(t, err)
		assert.Equal(t, expectedWriter.ID(), writer.ID())
		assert.Equal(t, expectedWriter.Name(), writer.Name())
	})

	t.Run("not found", func(t *testing.T) {
		t.Parallel()
		writerRepo := &mockWriterRepository{
			getByID: func(uint64) (*domain.Writer, error) {
				return nil, errors.New("not found")
			},
		}
		workRepo := &mockWorkRepository{}
		svc := service.NewWriterService(writerRepo, workRepo)

		_, err := svc.GetWriter(999)
		require.Error(t, err)
	})
}

func TestWriterService_ListWriters(t *testing.T) {
	t.Parallel()
	expectedWriters := []*domain.Writer{
		domain.NewWriter(1, "Jane Austen", 1775, nil, nil),
		domain.NewWriter(2, "Charles Dickens", 1812, nil, nil),
	}
	writerRepo := &mockWriterRepository{
		list: func(limit, offset int) ([]*domain.Writer, error) {
			return expectedWriters, nil
		},
	}
	workRepo := &mockWorkRepository{}
	svc := service.NewWriterService(writerRepo, workRepo)

	writers, err := svc.ListWriters(10, 0)
	require.NoError(t, err)
	assert.Len(t, writers, 2)
}

func TestWriterService_UpdateWriter(t *testing.T) {
	t.Parallel()
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		writerRepo := &mockWriterRepository{
			update: func(writer *domain.Writer) error {
				return nil
			},
		}
		workRepo := &mockWorkRepository{}
		svc := service.NewWriterService(writerRepo, workRepo)

		bio := "English novelist"
		err := svc.UpdateWriter(1, "Jane Austen", 1775, nil, &bio)
		assert.NoError(t, err)
	})

	t.Run("empty name", func(t *testing.T) {
		t.Parallel()
		writerRepo := &mockWriterRepository{}
		workRepo := &mockWorkRepository{}
		svc := service.NewWriterService(writerRepo, workRepo)

		err := svc.UpdateWriter(1, "", 1775, nil, nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "name is required")
	})

	t.Run("invalid birth year", func(t *testing.T) {
		t.Parallel()
		writerRepo := &mockWriterRepository{}
		workRepo := &mockWorkRepository{}
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
		writerRepo := &mockWriterRepository{
			delete: func(uint64) error {
				return nil
			},
		}
		workRepo := &mockWorkRepository{
			getByAuthorID: func(uint64) ([]*domain.Work, error) {
				return []*domain.Work{}, nil
			},
		}
		svc := service.NewWriterService(writerRepo, workRepo)

		err := svc.DeleteWriter(1)
		assert.NoError(t, err)
	})

	t.Run("cannot delete writer with works", func(t *testing.T) {
		t.Parallel()
		writerRepo := &mockWriterRepository{}
		workRepo := &mockWorkRepository{
			getByAuthorID: func(uint64) ([]*domain.Work, error) {
				return []*domain.Work{
					domain.NewWork(1, "Pride and Prejudice", 1),
				}, nil
			},
		}
		svc := service.NewWriterService(writerRepo, workRepo)

		err := svc.DeleteWriter(1)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "cannot delete writer with existing works")
	})
}
