package gorm

import (
	"github.com/what-writers-like/backend/internal/domain"
	"github.com/what-writers-like/backend/internal/infrastructure/database"
	"github.com/what-writers-like/backend/internal/repository"
	"gorm.io/gorm"
)

type workRepository struct {
	db *gorm.DB
}

func NewWorkRepository(db *database.Database) repository.WorkRepository {
	return &workRepository{db: db.DB()}
}

func (r *workRepository) Create(work *domain.Work) error {
	model := &database.WorkModel{
		ID:       work.ID(),
		Title:    work.Title(),
		AuthorID: work.AuthorID(),
	}
	return r.db.Create(model).Error
}

func (r *workRepository) GetByID(id uint64) (*domain.Work, error) {
	var model database.WorkModel
	if err := r.db.First(&model, id).Error; err != nil {
		return nil, err
	}
	return domain.NewWork(model.ID, model.Title, model.AuthorID), nil
}

func (r *workRepository) GetByAuthorID(authorID uint64) ([]*domain.Work, error) {
	var models []database.WorkModel
	if err := r.db.Where("author_id = ?", authorID).Find(&models).Error; err != nil {
		return nil, err
	}
	works := make([]*domain.Work, len(models))
	for i, m := range models {
		works[i] = domain.NewWork(m.ID, m.Title, m.AuthorID)
	}
	return works, nil
}

func (r *workRepository) List(limit, offset int) ([]*domain.Work, error) {
	var models []database.WorkModel
	if err := r.db.Limit(limit).Offset(offset).Find(&models).Error; err != nil {
		return nil, err
	}
	works := make([]*domain.Work, len(models))
	for i, m := range models {
		works[i] = domain.NewWork(m.ID, m.Title, m.AuthorID)
	}
	return works, nil
}

func (r *workRepository) Update(work *domain.Work) error {
	model := &database.WorkModel{
		ID:       work.ID(),
		Title:    work.Title(),
		AuthorID: work.AuthorID(),
	}
	return r.db.Save(model).Error
}

func (r *workRepository) Delete(id uint64) error {
	return r.db.Delete(&database.WorkModel{}, id).Error
}
