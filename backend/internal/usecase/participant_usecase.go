package usecase

import (
	"io"
	"siuji-backend/internal/entity"
	"siuji-backend/internal/model"
	"siuji-backend/internal/model/converter"
	"siuji-backend/internal/repository"
	"siuji-backend/pkg/password"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/xuri/excelize/v2"
)

type ParticipantUseCase struct {
	Log                         *logrus.Logger
	Validate                    *validator.Validate
	ParticipantPeriodRepository repository.ParticipantPeriodRepository
	UserRepository              repository.UserRepository
	PeriodRepository            repository.PeriodRepository
}

func NewParticipantUseCase(
	log *logrus.Logger,
	validate *validator.Validate,
	participantPeriodRepository repository.ParticipantPeriodRepository,
	userRepository repository.UserRepository,
	periodRepository repository.PeriodRepository,
) *ParticipantUseCase {
	return &ParticipantUseCase{
		Log:                         log,
		Validate:                    validate,
		ParticipantPeriodRepository: participantPeriodRepository,
		UserRepository:              userRepository,
		PeriodRepository:            periodRepository,
	}
}

func (c *ParticipantUseCase) AddParticipant(periodPublicID string, request *model.AddParticipantRequest) (*model.ParticipantResponse, error) {
	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("invalid add participant request: %+v", err)
		return nil, fiber.NewError(fiber.StatusBadRequest, "invalid request")
	}

	period, err := c.PeriodRepository.FindByPublicID(periodPublicID)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusNotFound, "period not found")
	}

	user, err := c.UserRepository.FindByEmail(request.Email)
	if err != nil {
		hashed, err := password.Hash(request.NIM)
		if err != nil {
			c.Log.Errorf("failed to hash password: %+v", err)
			return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to add participant")
		}

		user = &entity.User{
			PublicID:   uuid.New(),
			Name:       request.Name,
			Email:      request.Email,
			Password:   hashed,
			Role:       "participant",
			NIM:        request.NIM,
			University: request.University,
		}
		if err := c.UserRepository.Create(user); err != nil {
			c.Log.Errorf("failed to create participant user: %+v", err)
			return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to add participant")
		}
	}

	exists, err := c.ParticipantPeriodRepository.ExistsByPeriodAndUser(period.ID, user.ID)
	if err != nil {
		c.Log.Errorf("failed to check existing participant: %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to add participant")
	}
	if exists {
		return nil, fiber.NewError(fiber.StatusConflict, "participant already registered in this period")
	}

	participantPeriod := &entity.ParticipantPeriod{
		PublicID: uuid.New(),
		UserID:   user.ID,
		PeriodID: period.ID,
		Status:   "registered",
	}

	if err := c.ParticipantPeriodRepository.Create(participantPeriod); err != nil {
		c.Log.Errorf("failed to assign participant to period: %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to add participant")
	}

	participantPeriod.User = *user

	return converter.ParticipantToDetailResponse(participantPeriod, period.PublicID.String()), nil
}

func (c *ParticipantUseCase) GetAllByPeriod(periodPublicID, filter, sort string, limit, offset int) ([]model.ParticipantResponse, int64, error) {
	period, err := c.PeriodRepository.FindByPublicID(periodPublicID)
	if err != nil {
		return nil, 0, fiber.NewError(fiber.StatusNotFound, "period not found")
	}

	list, total, err := c.ParticipantPeriodRepository.FindAllByPeriodIDPagination(period.ID, filter, sort, limit, offset)
	if err != nil {
		c.Log.Errorf("failed to fetch participants: %+v", err)
		return nil, 0, fiber.NewError(fiber.StatusInternalServerError, "failed to fetch participants")
	}

	responses := make([]model.ParticipantResponse, 0, len(list))
	for _, pp := range list {
		responses = append(responses, *converter.ParticipantToResponse(&pp))
	}
	return responses, total, nil
}

func (c *ParticipantUseCase) GetDetail(periodPublicID, userPublicID string) (*model.ParticipantResponse, error) {
	period, err := c.PeriodRepository.FindByPublicID(periodPublicID)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusNotFound, "period not found")
	}

	pp, err := c.ParticipantPeriodRepository.FindByPeriodAndUserPublicID(period.ID, userPublicID)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusNotFound, "participant not found in this period")
	}

	return converter.ParticipantToDetailResponse(pp, period.PublicID.String()), nil
}

func (c *ParticipantUseCase) Update(periodPublicID, userPublicID string, request *model.UpdateParticipantRequest) (*model.ParticipantResponse, error) {
	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("invalid update participant request: %+v", err)
		return nil, fiber.NewError(fiber.StatusBadRequest, "invalid request")
	}

	period, err := c.PeriodRepository.FindByPublicID(periodPublicID)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusNotFound, "period not found")
	}

	pp, err := c.ParticipantPeriodRepository.FindByPeriodAndUserPublicID(period.ID, userPublicID)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusNotFound, "participant not found in this period")
	}

	user := &pp.User
	if request.Name != "" {
		user.Name = request.Name
	}
	if request.University != "" {
		user.University = request.University
	}
	if request.NIM != "" && request.NIM != user.NIM {
		user.NIM = request.NIM
		hashed, err := password.Hash(request.NIM)
		if err != nil {
			c.Log.Errorf("failed to hash new password: %+v", err)
			return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to update participant")
		}
		user.Password = hashed
	}
	if err := c.UserRepository.Update(user); err != nil {
		c.Log.Errorf("failed to update participant user data: %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to update participant")
	}

	if request.Status != "" {
		pp.Status = request.Status
	}
	if request.Score != nil {
		pp.Score = request.Score
	}
	if err := c.ParticipantPeriodRepository.Update(pp); err != nil {
		c.Log.Errorf("failed to update participant period: %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to update participant")
	}

	return converter.ParticipantToDetailResponse(pp, period.PublicID.String()), nil
}

func (c *ParticipantUseCase) Remove(periodPublicID, userPublicID string) error {
	period, err := c.PeriodRepository.FindByPublicID(periodPublicID)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "period not found")
	}

	if err := c.ParticipantPeriodRepository.DeleteByPeriodAndUserPublicID(period.ID, userPublicID); err != nil {
		return fiber.NewError(fiber.StatusNotFound, "participant not found in this period")
	}
	return nil
}

func (c *ParticipantUseCase) ImportFromExcel(periodPublicID string, file io.Reader) (*model.ImportParticipantResponse, error) {
	period, err := c.PeriodRepository.FindByPublicID(periodPublicID)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusNotFound, "period not found")
	}
	xlsx, err := excelize.OpenReader(file)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, "failed to read excel file")
	}
	defer xlsx.Close()

	sheetName := xlsx.GetSheetName(0)
	rows, err := xlsx.GetRows(sheetName)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, "failed to read rows from excel")
	}

	var participantsToAssign []entity.ParticipantPeriod
	var totalImported, totalSkipped int
	var errorDetails []model.ImportErrorDetail
	seenUserIDs := make(map[uint]bool)

	for i, row := range rows {
		if i == 0 {
			continue
		}
		if len(row) < 4 {
			totalSkipped++
			errorDetails = append(errorDetails, model.ImportErrorDetail{
				Row: i + 1, Email: "N/A",
				Message: "incomplete row columns (expected: Name, Email, NIM, University).",
			})
			continue
		}
		name := strings.TrimSpace(row[0])
		email := strings.TrimSpace(row[1])
		nim := strings.TrimSpace(row[2])
		university := strings.TrimSpace(row[3])

		if email == "" || !strings.Contains(email, "@") {
			totalSkipped++
			errorDetails = append(errorDetails, model.ImportErrorDetail{
				Row: i + 1, Email: email, Message: "invalid or empty email format.",
			})
			continue
		}

		user, err := c.UserRepository.FindByEmail(email)
		if err != nil {
			hashed, err := password.Hash(nim)
			if err != nil {
				totalSkipped++
				errorDetails = append(errorDetails, model.ImportErrorDetail{
					Row: i + 1, Email: email, Message: "failed to hash password for new user.",
				})
				continue
			}

			user = &entity.User{
				PublicID: uuid.New(),
				Name: name,
				Email: email,
				Password: hashed,
				Role: "participant",
				NIM: nim,
				University: university,
			}
			if err := c.UserRepository.Create(user); err != nil {
				totalSkipped++
				errorDetails = append(errorDetails, model.ImportErrorDetail{
					Row: i + 1, Email: email, Message: "Failed to create user in database.",
				})
				continue
			}
		}

		if seenUserIDs[user.ID] {
			totalSkipped++
			errorDetails = append(errorDetails, model.ImportErrorDetail{
				Row: i + 1, Email: email, Message: "Duplicate email within this Excel file.",
			})
			continue
		}

		exists, err := c.ParticipantPeriodRepository.ExistsByPeriodAndUser(period.ID, user.ID)
		if err != nil {
			totalSkipped++
			errorDetails = append(errorDetails, model.ImportErrorDetail{
				Row: i + 1, Email: email, Message: "Failed to check existing participant.",
			})
			continue
		}
		if exists {
			totalSkipped++
			errorDetails = append(errorDetails, model.ImportErrorDetail{
				Row: i + 1, Email: email, Message: "Participant already registered in this period.",
			})
			continue
		}

		participantsToAssign = append(participantsToAssign, entity.ParticipantPeriod{
			PublicID: uuid.New(),
			UserID:   user.ID,
			PeriodID: period.ID,
			Status:   "registered",
		})
		seenUserIDs[user.ID] = true
		totalImported++
	}

	if len(participantsToAssign) > 0 {
		if err := c.ParticipantPeriodRepository.BulkCreate(participantsToAssign); err != nil {
			c.Log.Errorf("failed to bulk assign participants: %+v", err)
			return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to bulk assign participants to period")
		}
	}

	return &model.ImportParticipantResponse{
		TotalImported: totalImported,
		TotalSkipped:  totalSkipped,
		Errors:        errorDetails,
	}, nil
}