package services

import (
	"errors"
	"mime/multipart"
	"siuji-backend/models"
	"siuji-backend/repositories"
	"siuji-backend/utils"
	"strings"

	"github.com/google/uuid"
	"github.com/xuri/excelize/v2"
)

type ParticipantService interface {
	AddParticipantToPeriod(periodPublicID string, req *models.ParticipantRequest) (*models.ParticipantPeriod, error)
	ImportParticipantsFromExcel(periodPublicID string, fileHeader *multipart.FileHeader) (*models.ImportParticipantResponse, error)
	GetParticipantsByPeriod(periodPublicID string, filter, sort string, limit, offset int) ([]models.ParticipantListResponse, int64, error)
	GetParticipantDetail(periodPublicID, userPublicID string) (*models.ParticipantResponse, error)
	UpdateParticipant(periodPublicID, userPublicID string, req *models.UpdateParticipantRequest) (*models.ParticipantResponse, error)
	RemoveParticipantFromPeriod(periodPublicID, userPublicID string) error
}

type participantService struct {
	participantRepo repositories.ParticipantRepository
	userRepo repositories.UserRepository
	periodRepo repositories.PeriodRepository
}

func NewParticipantService(
	participantRepo repositories.ParticipantRepository,
	userRepo repositories.UserRepository,
	periodRepo repositories.PeriodRepository,
) ParticipantService {
	return&participantService{
		participantRepo: participantRepo,
		userRepo: userRepo,
		periodRepo: periodRepo,
	}
}

func (s *participantService) AddParticipantToPeriod(periodPublicID string, req *models.ParticipantRequest) (*models.ParticipantPeriod, error) {
	period, err := s.periodRepo.FindByPublicID(periodPublicID)
	if err != nil {
		return nil, errors.New("period not found")
	}
	// cek email
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		hashedPassword, err := utils.HashPassword(req.NIM)
		if err != nil {
			return nil, errors.New("failed to hash password")
		}
		user = &models.User{
			PublicID: uuid.New(),
			Name: req.Name,
			Email: req.Email,
			Password: string(hashedPassword),
			Role: "participant",
			NIM: req.NIM,
			University: req.University,
		}
		if err := s.userRepo.Create(user); err != nil {
			return nil, err
		}
	}
	// cek user sudah ada di period belum
	existing, _ := s.participantRepo.FindByPeriodAndUserPublicID(period.ID, user.PublicID.String())
	if existing != nil {
		return nil, errors.New("participant already registered in this period")
	}
	//  Assign ke participant_periods
	participantPeriod := &models.ParticipantPeriod{
		PublicID: uuid.New().String(),
		UserID:   user.ID,
		PeriodID: period.ID,
		Status:   "registered",
	}
	if err := s.participantRepo.AssignParticipant(participantPeriod); err != nil {
		return nil, err
	}
	// Pasang objek user agar terbawa utuh ke response
	participantPeriod.User = *user
	return participantPeriod, nil
}

func (s *participantService) ImportParticipantsFromExcel(periodPublicID string, fileHeader *multipart.FileHeader) (*models.ImportParticipantResponse, error) {
	period, err := s.periodRepo.FindByPublicID(periodPublicID)
	if err != nil {
		return nil, errors.New("period not found")
	}
	// Buka file yang di-upload
	file, err := fileHeader.Open()
	if err != nil {
		return nil, errors.New("failed to open uploaded file")
	}
	defer file.Close()

	xlsx, err := excelize.OpenReader(file)
	if err != nil {
		return nil, errors.New("failed to read excel file")
	}
	defer xlsx.Close()

	sheetName := xlsx.GetSheetName(0)
	rows, err := xlsx.GetRows(sheetName)
	if err != nil {
		return nil, errors.New("failed to read rows from excel")
	}
	var participantsToAssign []models.ParticipantPeriod
	var totalImported int
	var totalSkipped int
	var errorDetails []models.ImportErrorDetail
	// Looping baris excel (asumsi baris 1 adalah header: Name, Email, NIM, University)
	for i, row := range rows {
		if i == 0 {
			continue // Lewati header
		}
		// Pastikan kolom minimal ada 4
		if len(row) < 4 {
			totalSkipped++
			errorDetails = append(errorDetails, models.ImportErrorDetail{
				Row:     i + 1,
				Email:   "N/A",
				Message: "Incomplete row columns (expected: Name, Email, NIM, University).",
			})
			continue
		}
		name := strings.TrimSpace(row[0])
		email := strings.TrimSpace(row[1])
		nim := strings.TrimSpace(row[2])
		university := strings.TrimSpace(row[3])
		// Validasi email sederhana
		if email == "" || !strings.Contains(email, "@") {
			totalSkipped++
			errorDetails = append(errorDetails, models.ImportErrorDetail{
				Row:     i + 1,
				Email:   email,
				Message: "Invalid or empty email format.",
			})
			continue
		}
		// Cek apakah User sudah ada berdasarkan Email
		user, err := s.userRepo.FindByEmail(email)
		if err != nil {
			// Jika belum ada, buat user baru dengan password default (NIM)
			hashedPassword, err := utils.HashPassword(nim)
			if err != nil {
				totalSkipped++
				errorDetails = append(errorDetails, models.ImportErrorDetail{
					Row:     i + 1,
					Email:   email,
					Message: "Failed to hash password for new user.",
				})
				continue
			}
			user = &models.User{
				PublicID:   uuid.New(),
				Name:       name,
				Email:      email,
				Password:   string(hashedPassword),
				Role:       "participant",
				NIM:        nim,
				University: university,
			}
			if err := s.userRepo.Create(user); err != nil {
				totalSkipped++
				errorDetails = append(errorDetails, models.ImportErrorDetail{
					Row:     i + 1,
					Email:   email,
					Message: "Failed to create user in database.",
				})
				continue
			}
		}
		// Cek apakah user sudah terdaftar di periode ini sebelumnya
		existing, _ := s.participantRepo.FindByPeriodAndUserPublicID(period.ID, user.PublicID.String())
		if existing != nil {
			totalSkipped++
			errorDetails = append(errorDetails, models.ImportErrorDetail{
				Row:     i + 1,
				Email:   email,
				Message: "Participant already registered in this period.",
			})
			continue
		}
		// Masukkan ke dalam list yang akan di-bulk insert
		participantsToAssign = append(participantsToAssign, models.ParticipantPeriod{
			PublicID: uuid.New().String(),
			UserID:   user.ID,
			PeriodID: period.ID,
			Status:   "registered",
		})
		totalImported++
	}
	// Jika ada data valid yang terkumpul, simpan sekaligus menggunakan BulkAssignParticipants
	if len(participantsToAssign) > 0 {
		if err := s.participantRepo.BulkAssignParticipants(participantsToAssign); err != nil {
			return nil, errors.New("failed to bulk assign participants to period")
		}
	}
	return &models.ImportParticipantResponse{
		TotalImported: totalImported,
		TotalSkipped:  totalSkipped,
		Errors:        errorDetails,
	}, nil
}

func (s *participantService) GetParticipantsByPeriod(periodPublicID string, filter, sort string, limit, offset int) ([]models.ParticipantListResponse, int64, error) {
	period, err := s.periodRepo.FindByPublicID(periodPublicID)
	if err != nil {
		return nil, 0, errors.New("period not found")
	}
	// Ambil data dari repository (menggunakan pagination yang sudah Anda buat)
	participants, total, err := s.participantRepo.FindByPeriodIDPagination(period.ID, filter, sort, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	// Mapping data ke bentuk Response DTO
	var responseList []models.ParticipantListResponse
	for _, p := range participants {
		responseList = append(responseList, models.ParticipantListResponse{
			PublicID: p.PublicID,
			User: models.UserResponse{
				PublicID:   p.User.PublicID.String(),
				Name:       p.User.Name,
				Email:      p.User.Email,
				NIM:        p.User.NIM,
				University: p.User.University,
				Role:       p.User.Role,
			},
			Status:    p.Status,
			Score:     p.Score,
			CreatedAt: p.CreatedAt,
		})
	}
	return responseList, total, nil
}

func (s *participantService) GetParticipantDetail(periodPublicID, userPublicID string) (*models.ParticipantResponse, error) {
	// Validasi apakah periode ada
	period, err := s.periodRepo.FindByPublicID(periodPublicID)
	if err != nil {
		return nil, errors.New("period not found")
	}
	// Ambil data relasi participant berdasarkan period ID dan user public ID
	participantPeriod, err := s.participantRepo.FindByPeriodAndUserPublicID(period.ID, userPublicID)
	if err != nil {
		return nil, err
	}
	// cek jika data tidak ditemukan
	if participantPeriod == nil {
		return nil, errors.New("participant not found in this period")
	}
	// Mapping data ke bentuk DTO Response
	response := &models.ParticipantResponse{
		PublicID:       participantPeriod.PublicID,
		PeriodPublicID: period.PublicID,
		User: models.UserResponse{
			PublicID:   participantPeriod.User.PublicID.String(),
			Name:       participantPeriod.User.Name,
			Email:      participantPeriod.User.Email,
			NIM:        participantPeriod.User.NIM,
			University: participantPeriod.User.University,
			Role:       participantPeriod.User.Role,
		},
		Status:    participantPeriod.Status,
		Score:     participantPeriod.Score,
		CreatedAt: participantPeriod.CreatedAt,
		UpdatedAt: participantPeriod.UpdatedAt,
	}
	return response, nil
}

func (s *participantService) UpdateParticipant(periodPublicID, userPublicID string, req *models.UpdateParticipantRequest) (*models.ParticipantResponse, error) {
	period, err := s.periodRepo.FindByPublicID(periodPublicID)
	if err != nil {
		return nil, errors.New("period not found")
	}
	// Ambil data relasi participant yang ada saat ini
	participantPeriod, err := s.participantRepo.FindByPeriodAndUserPublicID(period.ID, userPublicID)
	if err != nil {
		return nil, err
	}
	if participantPeriod == nil {
		return nil, errors.New("participant not found in this period")
	}
	// Update data User yang terikat
	user := &participantPeriod.User
	// Jika ada perubahan Nama, Email, Universitas, atau NIM
	if req.Name != "" {
		user.Name = req.Name
	}
	if req.University != "" {
		user.University = req.University
	}
	// Cek jika NIM di-update, maka Password juga ikut di-update
	if req.NIM != "" && req.NIM != user.NIM {
		user.NIM = req.NIM

		// Hash ulang password berdasarkan NIM baru
		hashedPassword, err := utils.HashPassword(req.NIM)
		if err != nil {
			return nil, errors.New("failed to hash new password")
		}
		user.Password = string(hashedPassword)
	}
	// Simpan perubahan data User ke database
	if err := s.userRepo.Update(user); err != nil {
		return nil, errors.New("failed to update user data")
	}
	// Update data di tabel participant_periods (misal status atau score jika ada di request)
	if req.Status != "" {
		if req.Status != utils.ParticipantStatusRegistered &&
			req.Status != utils.ParticipantStatusStarted &&
			req.Status != utils.ParticipantStatusCompleted {
				return nil, errors.New("invalid participant status.")
			}
		participantPeriod.Status = req.Status
	}
	if req.Score != nil {
		participantPeriod.Score = req.Score
	}
	if err := s.participantRepo.Update(participantPeriod); err != nil {
		return nil, errors.New("failed to update participant period")
	}
	// Mapping ke Response DTO
	response := &models.ParticipantResponse{
		PublicID:       participantPeriod.PublicID,
		PeriodPublicID: period.PublicID,
		User: models.UserResponse{
			PublicID:   user.PublicID.String(),
			Name:       user.Name,
			Email:      user.Email,
			NIM:        user.NIM,
			University: user.University,
			Role:       user.Role,
		},
		Status:    participantPeriod.Status,
		Score:     participantPeriod.Score,
		CreatedAt: participantPeriod.CreatedAt,
		UpdatedAt: participantPeriod.UpdatedAt,
	}

	return response, nil
}

func (s *participantService) RemoveParticipantFromPeriod(periodPublicID, userPublicID string) error {
	period, err := s.periodRepo.FindByPublicID(periodPublicID)
	if err != nil {
		return errors.New("period not found")
	}
	// Cari data user berdasarkan userPublicID untuk mendapatkan internal ID (uint)
	user, err := s.userRepo.FindByPublicID(userPublicID) // Sesuaikan dengan nama method di userRepo Anda
	if err != nil || user == nil {
		return errors.New("user not found")
	}
	// Panggil repository participant yang sudah ada
	err = s.participantRepo.RemoveFromPeriod(period.ID, user.ID)
	if err != nil {
		return err
	}
	return nil
}