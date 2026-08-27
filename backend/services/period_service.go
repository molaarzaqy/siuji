package services

import (
	"errors"
	"siuji-backend/models"
	"siuji-backend/repositories"
	"siuji-backend/utils"

	"github.com/google/uuid"
)

type PeriodService interface{
	CreatePeriod(req *models.PeriodRequest) (*models.Period, error)
	GetAllPeriods(filter, sort string, limit, offset int) ([]models.Period, int64, error)
	GetPeriodByPublicID(publicID string) (*models.PeriodDetailResponse, error)
	UpdatePeriod(publicID string, req *models.PeriodRequest) (*models.Period, error)
	DeletePeriod(publicID string) error
	
	// Manajemen Section dalam Period
	AddSectionToPeriod(periodPublicID, sectionPublicID string, position int) (*models.PeriodSection, error)
	RemoveSectionFromPeriod(periodPublicID, sectionPublicID string) error
	ReorderSections(periodPublicID string, sectionPublicIDs []string) error
}

type periodService struct {
	periodRepo repositories.PeriodRepository
	sectionRepo repositories.SectionRepository
}

func NewPeriodService(periodRepo repositories.PeriodRepository,sectionRepo repositories.SectionRepository) PeriodService {
	return &periodService{
		periodRepo: periodRepo,
		sectionRepo: sectionRepo,
	}
}

func (s *periodService) CreatePeriod(req *models.PeriodRequest) (*models.Period, error) {
	period := &models.Period{
		PublicID:            uuid.New().String(),
		Title:               req.Title,
		Month:               req.Month,
		Year:                req.Year,
		Status:              req.Status,
		CertificateURL:      req.CertificateURL,
		CertificateExpMonth: req.CertificateExpMonth,
		MinPassingGrade:     req.MinPassingGrade,
		MaxPassingGrade:     req.MaxPassingGrade,
		StartTime:           req.StartTime,
		EndTime:             req.EndTime,
	}

	err := s.periodRepo.Create(period)
	if err != nil {
		return nil, err
	}
	return period, nil
}

func (s *periodService) GetAllPeriods(filter, sort string, limit, offset int) ([]models.Period, int64, error) {
	return s.periodRepo.FindAllPagination(filter, sort, limit, offset)
}

func (s *periodService) GetPeriodByPublicID(publicID string) (*models.PeriodDetailResponse, error) {
	period, err := s.periodRepo.FindByPublicID(publicID)
    if err != nil {
        return nil, err
    }
    if period == nil {
        return nil, errors.New("period not found")
    }

    // ambil list section
    var sections []models.PeriodSectionResponse
    for _, ps := range period.PeriodSections {
        sections = append(sections, models.PeriodSectionResponse{
            PeriodSectionPublicID: ps.PublicID,
            SectionPublicID: ps.Section.PublicID,
            Title: ps.Section.Title,
            Position: ps.Position,
        })
    }
    // Buat response struct
    response := &models.PeriodDetailResponse{
        PublicID: period.PublicID,
        Title: period.Title,
        Month: period.Month,
        Year: period.Year,
        DueDate: period.DueDate,
        Status: period.Status,
        CertificateURL: period.CertificateURL,
        CertificateExpMonth: period.CertificateExpMonth,
        MinPassingGrade: period.MinPassingGrade,
        MaxPassingGrade: period.MaxPassingGrade,
        StartTime: period.StartTime,
        EndTime: period.EndTime,
        Sections: sections,
        CreatedAt: period.CreatedAt,
        UpdatedAt: period.UpdatedAt,
    }

    return response, nil
}

func (s *periodService) UpdatePeriod(publicID string, req *models.PeriodRequest) (*models.Period, error) {
	period, err := s.periodRepo.FindByPublicID(publicID)
	if err != nil {
		return nil, err
	}
	if period == nil {
		return nil, errors.New("period not found")
	}
	if req.Status != "" {
		if req.Status != utils.PeriodStatusDraft &&
		   req.Status != utils.PeriodStatusPublished &&
		   req.Status != utils.PeriodStatusClosed {
			return nil, errors.New("invalid period status")
		}
		period.Status = req.Status
	}

	period.Title = req.Title
	period.Month = req.Month
	period.Year = req.Year
	period.DueDate = req.DueDate
	period.CertificateURL = req.CertificateURL
	period.CertificateExpMonth = req.CertificateExpMonth
	period.MinPassingGrade = req.MinPassingGrade
	period.MaxPassingGrade = req.MaxPassingGrade
	period.StartTime = req.StartTime
	period.EndTime = req.EndTime

	err = s.periodRepo.Update(period)
	if err != nil {
		return nil, err
	}
	return period, nil
}

func (s *periodService) DeletePeriod(publicID string) error {
	period, err := s.periodRepo.FindByPublicID(publicID)
	if err != nil {
		return  err
	}
	if period == nil {
		return errors.New("period not found")
	}

	return s.periodRepo.Delete(publicID)
}

func (s *periodService) AddSectionToPeriod(periodPublicID, sectionPublicID string, position int) (*models.PeriodSection, error) {
	period, err := s.periodRepo.FindByPublicID(periodPublicID)
	if err != nil || period == nil {
		return nil, errors.New("period not found")
	}

	section, err := s.sectionRepo.FindByPublicID(sectionPublicID)
	if err != nil || section == nil {
		return nil, errors.New("section not found")
	}

	if position <= 0 {
		maxPos, err := s.periodRepo.GetMaxPositionInPeriod(period.ID)
		if err != nil {
			return nil, err
		}
		position = maxPos + 1
	}

	periodSection := &models.PeriodSection{
		PublicID:  uuid.New().String(),
		PeriodID:  period.ID,
		SectionID: section.ID,
		Position:  position,
	}

	err = s.periodRepo.AddSectionToPeriod(periodSection)
	if err != nil {
		return nil, err
	}

	// Load relasi Section agar response service lengkap (termasuk title section)
	periodSection.Section = *section

	return periodSection, nil
}

func (s *periodService) RemoveSectionFromPeriod(periodPublicID, sectionPublicID string) error {
	period, err := s.periodRepo.FindByPublicID(periodPublicID)
	if err != nil || period == nil {
		return errors.New("period not found")
	}

	section, err := s.sectionRepo.FindByPublicID(sectionPublicID)
	if err != nil || section == nil {
		return errors.New("section not found")
	}

	return s.periodRepo.RemoveSectionFromPeriod(period.ID, section.ID)
}

func (s *periodService) ReorderSections(periodPublicID string, sectionPublicIDs []string) error {
period, err := s.periodRepo.FindByPublicID(periodPublicID)
	if err != nil || period == nil {
		return errors.New("period not found")
	}

	var updates []repositories.SectionPositionUpdate
	for index, secPubID := range sectionPublicIDs {
		section, err := s.sectionRepo.FindByPublicID(secPubID)
		if err != nil || section == nil {
			return errors.New("one or more sections not found")
		}
		updates = append(updates, repositories.SectionPositionUpdate{
			PeriodID:  period.ID,
			SectionID: section.ID,
			Position:  index + 1,
		})
	}

	return s.periodRepo.UpdateSectionPositions(updates)
}