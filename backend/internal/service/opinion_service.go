package service

import (
	"errors"

	"github.com/what-writers-like/backend/internal/domain"
	"github.com/what-writers-like/backend/internal/repository"
)

type OpinionService interface {
	CreateOpinion(
		writerID, workID uint64,
		sentiment bool,
		quote, source string,
		page *string,
		statementYear *int,
	) (*domain.Opinion, error)
	GetOpinionsByWriter(writerID uint64) ([]*domain.Opinion, error)
	GetOpinionsByWork(workID uint64) ([]*domain.Opinion, error)
	GetOpinion(writerID, workID uint64) (*domain.Opinion, error)
	ListOpinions(limit, offset int) ([]*domain.Opinion, error)
	UpdateOpinion(writerID, workID uint64, sentiment bool, quote, source string, page *string, statementYear *int) error
	DeleteOpinion(writerID, workID uint64) error
}

type opinionService struct {
	opinionRepo repository.OpinionRepository
	writerRepo  repository.WriterRepository
	workRepo    repository.WorkRepository
}

func NewOpinionService(
	opinionRepo repository.OpinionRepository,
	writerRepo repository.WriterRepository,
	workRepo repository.WorkRepository,
) OpinionService {
	return &opinionService{
		opinionRepo: opinionRepo,
		writerRepo:  writerRepo,
		workRepo:    workRepo,
	}
}

func (s *opinionService) CreateOpinion(
	writerID, workID uint64,
	sentiment bool,
	quote, source string,
	page *string,
	statementYear *int,
) (*domain.Opinion, error) {
	if quote == "" {
		return nil, errors.New("quote is required")
	}
	if source == "" {
		return nil, errors.New("source is required")
	}

	work, err := s.workRepo.GetByID(workID)
	if err != nil {
		return nil, errors.New("work not found")
	}

	if work.AuthorID() == writerID {
		return nil, errors.New("writer cannot express opinion about their own work")
	}

	_, err = s.writerRepo.GetByID(writerID)
	if err != nil {
		return nil, errors.New("writer not found")
	}

	opinion := domain.NewOpinion(writerID, workID, sentiment, quote, source, page, statementYear)
	if err := s.opinionRepo.Create(opinion); err != nil {
		return nil, err
	}
	return opinion, nil
}

func (s *opinionService) GetOpinionsByWriter(writerID uint64) ([]*domain.Opinion, error) {
	return s.opinionRepo.GetByWriterID(writerID)
}

func (s *opinionService) GetOpinionsByWork(workID uint64) ([]*domain.Opinion, error) {
	return s.opinionRepo.GetByWorkID(workID)
}

func (s *opinionService) GetOpinion(writerID, workID uint64) (*domain.Opinion, error) {
	return s.opinionRepo.GetByWriterAndWork(writerID, workID)
}

func (s *opinionService) ListOpinions(limit, offset int) ([]*domain.Opinion, error) {
	return s.opinionRepo.List(limit, offset)
}

func (s *opinionService) UpdateOpinion(
	writerID, workID uint64,
	sentiment bool,
	quote, source string,
	page *string,
	statementYear *int,
) error {
	if quote == "" {
		return errors.New("quote is required")
	}
	if source == "" {
		return errors.New("source is required")
	}

	work, err := s.workRepo.GetByID(workID)
	if err != nil {
		return errors.New("work not found")
	}

	if work.AuthorID() == writerID {
		return errors.New("writer cannot express opinion about their own work")
	}

	opinion := domain.NewOpinion(writerID, workID, sentiment, quote, source, page, statementYear)
	return s.opinionRepo.Update(opinion)
}

func (s *opinionService) DeleteOpinion(writerID, workID uint64) error {
	return s.opinionRepo.Delete(writerID, workID)
}
