package repository

import "github.com/what-writers-like/backend/internal/domain"

type WorkRepository interface {
	Create(work *domain.Work) error
	GetByID(id uint64) (*domain.Work, error)
	GetByAuthorID(authorID uint64) ([]*domain.Work, error)
	List(limit, offset int) ([]*domain.Work, error)
	Search(query string, limit, offset int) ([]*domain.Work, error)
	Update(work *domain.Work) error
	Delete(id uint64) error
}
