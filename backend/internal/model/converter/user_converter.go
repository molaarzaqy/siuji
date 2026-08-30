package converter

import (
	"siuji-backend/internal/entity"
	"siuji-backend/internal/model"
)

func UserToResponse(user *entity.User) *model.UserResponse {
	return &model.UserResponse{
		PublicID: user.PublicID.String(),
		Name: user.Name,
		Email: user.Email,
		Role: user.Role,
		NIM: user.NIM,
		University: user.University,
		UpdatedAt: user.UpdatedAt,
	}
}