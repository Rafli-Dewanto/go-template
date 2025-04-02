package service

import (
	"context"
	"errors"
	"math"
	"time"

	"github.com/Rafli-Dewanto/go-template/internal/entity"
	"github.com/Rafli-Dewanto/go-template/internal/model"
	"github.com/Rafli-Dewanto/go-template/internal/model/converter"
	"github.com/Rafli-Dewanto/go-template/internal/repository"
	"github.com/Rafli-Dewanto/go-template/internal/utils"
)

var (
	ErrUserNotFound         = errors.New("user not found")
	ErrInvalidInput         = errors.New("invalid input")
	ErrUserAlreadyExists    = errors.New("user already exists")
	ErrRequestTimeout       = errors.New("request timeout")
	ErrUsernameAlreadyTaken = errors.New("username already taken")
)

type UserService interface {
	Create(ctx context.Context, user *model.CreateUserRequest) error
	GetByID(ctx context.Context, id int64) (*model.UserResponse, error)
	List(ctx context.Context, query *model.PaginationQuery) (*model.Response, error)
	Update(ctx context.Context, user model.UpdateUserRequest) error
	SoftDelete(ctx context.Context, id int64) error
}

type userService struct {
	repo   repository.UserRepository
	logger *utils.Logger
}

func NewUserService(repo repository.UserRepository, logger *utils.Logger) UserService {
	return &userService{repo: repo, logger: logger}
}

func (s *userService) Create(ctx context.Context, user *model.CreateUserRequest) error {
	if ctx.Err() != nil {
		s.logger.Warning("Request timeout: operation took longer than 10 seconds")
		return ErrRequestTimeout
	}

	if user.Username == "" || user.Email == "" {
		s.logger.Warning("Invalid input for user creation: %v", ErrInvalidInput)
		return ErrInvalidInput
	}

	// Check if user already exists
	existingUser, err := s.repo.GetByEmailOrUsername(ctx, user.Email, user.Username)
	if err == context.DeadlineExceeded || err == context.Canceled {
		s.logger.Warning("Database query timed out")
		return ErrRequestTimeout
	}

	if err == nil && existingUser != nil {
		s.logger.Warning("User with email or username already exists")
		return ErrUserAlreadyExists
	}

	newUser := &entity.User{
		Username:  user.Username,
		Email:     user.Email,
		Password:  user.Password,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Insert user into DB
	err = s.repo.Create(ctx, newUser)
	if err == context.DeadlineExceeded || err == context.Canceled {
		s.logger.Warning("Database insert timed out")
		return ErrRequestTimeout
	}
	if err != nil {
		return err
	}

	return nil
}

func (s *userService) GetByID(ctx context.Context, id int64) (*model.UserResponse, error) {
	if id <= 0 {
		s.logger.Warning("Invalid input for user retrieval: %v", ErrInvalidInput)
		return nil, ErrInvalidInput
	}

	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Warning("User not found: %v", ErrUserNotFound)
		return nil, ErrUserNotFound
	}
	return converter.ToUserResponse(user), nil
}

func (s *userService) List(ctx context.Context, query *model.PaginationQuery) (*model.Response, error) {
	if query.Limit <= 0 {
		query.Limit = 10
	}

	users, total, err := s.repo.List(ctx, query)
	if err != nil {
		s.logger.Warning("Failed to list users: %v", err)
		return nil, err
	}

	totalPages := int(math.Ceil(float64(total) / float64(query.Limit)))
	userResponses := make([]*model.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = converter.ToUserResponse(user)
	}

	return &model.Response{
		Message: "Users retrieved successfully",
		Data:    userResponses,
		Meta: &model.PaginatedMeta{
			Total:       total,
			CurrentPage: int64(query.Page),
			PerPage:     int64(query.Limit),
			LastPage:    totalPages,
			HasNextPage: int64(query.Page) < int64(totalPages),
			HasPrevPage: int64(query.Page) > 1,
		},
	}, nil
}

func (s *userService) Update(ctx context.Context, user model.UpdateUserRequest) error {
	if user.ID <= 0 {
		s.logger.Warning("Invalid input for user update: %v", ErrInvalidInput)
		return ErrInvalidInput
	}

	existingUser, err := s.repo.GetByID(ctx, user.ID)
	if err != nil {
		s.logger.Warning("User not found: %v", ErrUserNotFound)
		return err
	}

	updatedUser := &entity.User{
		ID:       user.ID,
		Username: existingUser.Username,
		Email:    existingUser.Email,
	}

	// Update username if provided
	if user.Username != nil {
		// Check if new username is already taken
		_, err := s.repo.GetByUsername(ctx, *user.Username)
		if err == nil {
			s.logger.Warning("Username is already taken: %v", ErrUsernameAlreadyTaken)
			return ErrUsernameAlreadyTaken
		}
		updatedUser.Username = *user.Username
	}

	// Update email if provided
	if user.Email != nil {
		updatedUser.Email = *user.Email
	}

	err = s.repo.Update(ctx, updatedUser)
	if err != nil {
		return err
	}
	return nil
}

func (s *userService) SoftDelete(ctx context.Context, id int64) error {
	if id <= 0 {
		return ErrInvalidInput
	}

	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Warning("User not found: %v", ErrUserNotFound)
		return err
	}

	return s.repo.SoftDelete(ctx, id)
}
