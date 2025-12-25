package repository

import "github.com/what-writers-like/backend/internal/domain"

type WriterRepository interface {
	Create(writer *domain.Writer) error
	GetByID(id uint64) (*domain.Writer, error)
	List(limit, offset int) ([]*domain.Writer, error)
	Update(writer *domain.Writer) error
	Delete(id uint64) error
}
