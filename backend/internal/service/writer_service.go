package service

import (
	"errors"

	"github.com/what-writers-like/backend/internal/domain"
	"github.com/what-writers-like/backend/internal/repository"
)

type WriterService interface {
	CreateWriter(name string, birthYear int, deathYear *int, bio *string) (*domain.Writer, error)
	GetWriter(id uint64) (*domain.Writer, error)
	ListWriters(limit, offset int) ([]*domain.Writer, error)
	UpdateWriter(id uint64, name string, birthYear int, deathYear *int, bio *string) error
	DeleteWriter(id uint64) error
}

type writerService struct {
	writerRepo repository.WriterRepository
	workRepo   repository.WorkRepository
}

func NewWriterService(writerRepo repository.WriterRepository, workRepo repository.WorkRepository) WriterService {
	return &writerService{
		writerRepo: writerRepo,
		workRepo:   workRepo,
	}
}

func (s *writerService) CreateWriter(name string, birthYear int, deathYear *int, bio *string) (*domain.Writer, error) {
	if name == "" {
		return nil, errors.New("name is required")
	}
	if birthYear <= 0 {
		return nil, errors.New("birth year must be positive")
	}

	id := uint64(1)
	writers, _ := s.writerRepo.List(1, 0)
	if len(writers) > 0 {
		maxID := uint64(0)
		allWriters, _ := s.writerRepo.List(1000, 0)
		for _, w := range allWriters {
			if w.ID() > maxID {
				maxID = w.ID()
			}
		}
		id = maxID + 1
	}

	writer := domain.NewWriter(id, name, birthYear, deathYear, bio)
	if err := s.writerRepo.Create(writer); err != nil {
		return nil, err
	}
	return writer, nil
}

func (s *writerService) GetWriter(id uint64) (*domain.Writer, error) {
	return s.writerRepo.GetByID(id)
}

func (s *writerService) ListWriters(limit, offset int) ([]*domain.Writer, error) {
	return s.writerRepo.List(limit, offset)
}

func (s *writerService) UpdateWriter(id uint64, name string, birthYear int, deathYear *int, bio *string) error {
	if name == "" {
		return errors.New("name is required")
	}
	if birthYear <= 0 {
		return errors.New("birth year must be positive")
	}

	writer := domain.NewWriter(id, name, birthYear, deathYear, bio)
	return s.writerRepo.Update(writer)
}

func (s *writerService) DeleteWriter(id uint64) error {
	works, err := s.workRepo.GetByAuthorID(id)
	if err != nil {
		return err
	}
	if len(works) > 0 {
		return errors.New("cannot delete writer with existing works")
	}
	return s.writerRepo.Delete(id)
}
