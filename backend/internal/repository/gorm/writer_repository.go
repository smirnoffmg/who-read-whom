package gorm

import (
	"github.com/what-writers-like/backend/internal/domain"
	"github.com/what-writers-like/backend/internal/infrastructure/database"
	"github.com/what-writers-like/backend/internal/repository"
	"gorm.io/gorm"
)

type writerRepository struct {
	db *gorm.DB
}

func NewWriterRepository(db *database.Database) repository.WriterRepository {
	return &writerRepository{db: db.DB()}
}

func (r *writerRepository) Create(writer *domain.Writer) error {
	model := &database.WriterModel{
		ID:        writer.ID(),
		Name:      writer.Name(),
		BirthYear: writer.BirthYear(),
		DeathYear: writer.DeathYear(),
		Bio:       writer.Bio(),
	}
	return r.db.Create(model).Error
}

func (r *writerRepository) GetByID(id uint64) (*domain.Writer, error) {
	var model database.WriterModel
	if err := r.db.First(&model, id).Error; err != nil {
		return nil, err
	}
	return domain.NewWriter(model.ID, model.Name, model.BirthYear, model.DeathYear, model.Bio), nil
}

func (r *writerRepository) List(limit, offset int) ([]*domain.Writer, error) {
	var models []database.WriterModel
	if err := r.db.Limit(limit).Offset(offset).Find(&models).Error; err != nil {
		return nil, err
	}
	writers := make([]*domain.Writer, len(models))
	for i, m := range models {
		writers[i] = domain.NewWriter(m.ID, m.Name, m.BirthYear, m.DeathYear, m.Bio)
	}
	return writers, nil
}

func (r *writerRepository) Search(query string, limit, offset int) ([]*domain.Writer, error) {
	var models []database.WriterModel
	// Use PostgreSQL fuzzy search with similarity threshold of 0.3
	// similarity() function from pg_trgm returns a value between 0 and 1
	searchSQL := `
		SELECT * FROM writers 
		WHERE similarity(name, ?) > 0.3 
		   OR (bio IS NOT NULL AND similarity(bio, ?) > 0.3)
		ORDER BY 
			GREATEST(similarity(name, ?), COALESCE(similarity(bio, ?), 0)) DESC
		LIMIT ? OFFSET ?
	`
	if err := r.db.Raw(searchSQL, query, query, query, query, limit, offset).Scan(&models).Error; err != nil {
		return nil, err
	}
	writers := make([]*domain.Writer, len(models))
	for i, m := range models {
		writers[i] = domain.NewWriter(m.ID, m.Name, m.BirthYear, m.DeathYear, m.Bio)
	}
	return writers, nil
}

func (r *writerRepository) Update(writer *domain.Writer) error {
	model := &database.WriterModel{
		ID:        writer.ID(),
		Name:      writer.Name(),
		BirthYear: writer.BirthYear(),
		DeathYear: writer.DeathYear(),
		Bio:       writer.Bio(),
	}
	return r.db.Save(model).Error
}

func (r *writerRepository) Delete(id uint64) error {
	return r.db.Delete(&database.WriterModel{}, id).Error
}
