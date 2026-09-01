package converter

import (
	"siuji-backend/internal/entity"
	"siuji-backend/internal/model"
)

func PeriodToResponse(period *entity.Period) *model.PeriodResponse {
	return &model.PeriodResponse{
		PublicID:            period.PublicID.String(),
		Title:                period.Title,
		Month:                period.Month,
		Year:                 period.Year,
		Status:               period.Status,
		CertificateURL:       period.CertificateURL,
		CertificateExpMonth:  period.CertificateExpMonth,
		MinPassingGrade:      period.MinPassingGrade,
		MaxPassingGrade:      period.MaxPassingGrade,
		StartTime:            period.StartTime,
		EndTime:              period.EndTime,
		CreatedAt:            period.CreatedAt,
		UpdatedAt:            period.UpdatedAt,
	}
}

func PeriodToDetailResponse(period *entity.Period) *model.PeriodDetailResponse {
	sections := make([]model.PeriodSectionResponse, 0, len(period.PeriodSections))
	for _, ps := range period.PeriodSections {
		sections = append(sections, *PeriodSectionToResponse(&ps))
	}

	return &model.PeriodDetailResponse{
		PeriodResponse: *PeriodToResponse(period),
		Sections:       sections,
	}
}