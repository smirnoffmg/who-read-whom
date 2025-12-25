package gorm

import (
	"github.com/what-writers-like/backend/internal/domain"
	"github.com/what-writers-like/backend/internal/infrastructure/database"
	"github.com/what-writers-like/backend/internal/repository"
	"gorm.io/gorm"
)

type opinionRepository struct {
	db *gorm.DB
}

func NewOpinionRepository(db *database.Database) repository.OpinionRepository {
	return &opinionRepository{db: db.DB()}
}

func (r *opinionRepository) Create(opinion *domain.Opinion) error {
	model := &database.OpinionModel{
		WriterID:      opinion.WriterID(),
		WorkID:        opinion.WorkID(),
		Sentiment:     opinion.Sentiment(),
		Quote:         opinion.Quote(),
		Source:        opinion.Source(),
		Page:          opinion.Page(),
		StatementYear: opinion.StatementYear(),
	}
	return r.db.Create(model).Error
}

func (r *opinionRepository) GetByWriterID(writerID uint64) ([]*domain.Opinion, error) {
	var models []database.OpinionModel
	if err := r.db.Where("writer_id = ?", writerID).Find(&models).Error; err != nil {
		return nil, err
	}
	opinions := make([]*domain.Opinion, len(models))
	for i, m := range models {
		opinions[i] = domain.NewOpinion(m.WriterID, m.WorkID, m.Sentiment, m.Quote, m.Source, m.Page, m.StatementYear)
	}
	return opinions, nil
}

func (r *opinionRepository) GetByWorkID(workID uint64) ([]*domain.Opinion, error) {
	var models []database.OpinionModel
	if err := r.db.Where("work_id = ?", workID).Find(&models).Error; err != nil {
		return nil, err
	}
	opinions := make([]*domain.Opinion, len(models))
	for i, m := range models {
		opinions[i] = domain.NewOpinion(m.WriterID, m.WorkID, m.Sentiment, m.Quote, m.Source, m.Page, m.StatementYear)
	}
	return opinions, nil
}

func (r *opinionRepository) GetByWriterAndWork(writerID, workID uint64) (*domain.Opinion, error) {
	var model database.OpinionModel
	if err := r.db.Where("writer_id = ? AND work_id = ?", writerID, workID).First(&model).Error; err != nil {
		return nil, err
	}
	return domain.NewOpinion(
		model.WriterID,
		model.WorkID,
		model.Sentiment,
		model.Quote,
		model.Source,
		model.Page,
		model.StatementYear,
	), nil
}

func (r *opinionRepository) List(limit, offset int) ([]*domain.Opinion, error) {
	var models []database.OpinionModel
	if err := r.db.Limit(limit).Offset(offset).Find(&models).Error; err != nil {
		return nil, err
	}
	opinions := make([]*domain.Opinion, len(models))
	for i, m := range models {
		opinions[i] = domain.NewOpinion(m.WriterID, m.WorkID, m.Sentiment, m.Quote, m.Source, m.Page, m.StatementYear)
	}
	return opinions, nil
}

func (r *opinionRepository) Update(opinion *domain.Opinion) error {
	model := &database.OpinionModel{
		WriterID:      opinion.WriterID(),
		WorkID:        opinion.WorkID(),
		Sentiment:     opinion.Sentiment(),
		Quote:         opinion.Quote(),
		Source:        opinion.Source(),
		Page:          opinion.Page(),
		StatementYear: opinion.StatementYear(),
	}
	return r.db.Save(model).Error
}

func (r *opinionRepository) Delete(writerID, workID uint64) error {
	return r.db.Where("writer_id = ? AND work_id = ?", writerID, workID).Delete(&database.OpinionModel{}).Error
}
