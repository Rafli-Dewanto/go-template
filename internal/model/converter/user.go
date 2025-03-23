package converter

import (
	"github.com/Rafli-Dewanto/go-template/internal/entity"
	"github.com/Rafli-Dewanto/go-template/internal/model"
)

func ToUserResponse(user *entity.User) *model.UserResponse {
	return &model.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func ToUsersResponse(users []*entity.User, meta *model.PaginatedMeta) *model.Response {
	userResponses := make([]*model.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = ToUserResponse(user)
	}

	return &model.Response{
		Success: true,
		Message: "Users retrieved successfully",
		Data:    userResponses,
		Meta:    meta,
	}
}
