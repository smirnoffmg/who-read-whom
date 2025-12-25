package repository

import "github.com/what-writers-like/backend/internal/domain"

type OpinionRepository interface {
	Create(opinion *domain.Opinion) error
	GetByWriterID(writerID uint64) ([]*domain.Opinion, error)
	GetByWorkID(workID uint64) ([]*domain.Opinion, error)
	GetByWriterAndWork(writerID, workID uint64) (*domain.Opinion, error)
	List(limit, offset int) ([]*domain.Opinion, error)
	Update(opinion *domain.Opinion) error
	Delete(writerID, workID uint64) error
}
