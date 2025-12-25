package service

import (
	"errors"

	"github.com/what-writers-like/backend/internal/domain"
	"github.com/what-writers-like/backend/internal/repository"
)

type WorkService interface {
	CreateWork(title string, authorID uint64) (*domain.Work, error)
	GetWork(id uint64) (*domain.Work, error)
	GetWorksByAuthor(authorID uint64) ([]*domain.Work, error)
	ListWorks(limit, offset int) ([]*domain.Work, error)
	SearchWorks(query string, limit, offset int) ([]*domain.Work, error)
	UpdateWork(id uint64, title string, authorID uint64) error
	DeleteWork(id uint64) error
}

type workService struct {
	workRepo   repository.WorkRepository
	writerRepo repository.WriterRepository
}

func NewWorkService(workRepo repository.WorkRepository, writerRepo repository.WriterRepository) WorkService {
	return &workService{
		workRepo:   workRepo,
		writerRepo: writerRepo,
	}
}

func (s *workService) CreateWork(title string, authorID uint64) (*domain.Work, error) {
	if title == "" {
		return nil, errors.New("title is required")
	}

	_, err := s.writerRepo.GetByID(authorID)
	if err != nil {
		return nil, errors.New("author not found")
	}

	id := uint64(1)
	works, _ := s.workRepo.List(1, 0)
	if len(works) > 0 {
		maxID := uint64(0)
		allWorks, _ := s.workRepo.List(1000, 0)
		for _, w := range allWorks {
			if w.ID() > maxID {
				maxID = w.ID()
			}
		}
		id = maxID + 1
	}

	work := domain.NewWork(id, title, authorID)
	if err := s.workRepo.Create(work); err != nil {
		return nil, err
	}
	return work, nil
}

func (s *workService) GetWork(id uint64) (*domain.Work, error) {
	return s.workRepo.GetByID(id)
}

func (s *workService) GetWorksByAuthor(authorID uint64) ([]*domain.Work, error) {
	return s.workRepo.GetByAuthorID(authorID)
}

func (s *workService) ListWorks(limit, offset int) ([]*domain.Work, error) {
	return s.workRepo.List(limit, offset)
}

func (s *workService) SearchWorks(query string, limit, offset int) ([]*domain.Work, error) {
	if query == "" {
		return s.workRepo.List(limit, offset)
	}
	return s.workRepo.Search(query, limit, offset)
}

func (s *workService) UpdateWork(id uint64, title string, authorID uint64) error {
	if title == "" {
		return errors.New("title is required")
	}

	// Check if work exists
	_, err := s.workRepo.GetByID(id)
	if err != nil {
		return errors.New("work not found")
	}

	_, err = s.writerRepo.GetByID(authorID)
	if err != nil {
		return errors.New("author not found")
	}

	work := domain.NewWork(id, title, authorID)
	return s.workRepo.Update(work)
}

func (s *workService) DeleteWork(id uint64) error {
	// Check if work exists
	_, err := s.workRepo.GetByID(id)
	if err != nil {
		return errors.New("work not found")
	}
	return s.workRepo.Delete(id)
}
