package repository

import (
	"fmt"

	"github.com/Rafli-Dewanto/go-template/internal/entity"
	"github.com/Rafli-Dewanto/go-template/internal/model"
	"github.com/Rafli-Dewanto/go-template/internal/utils"
	"github.com/jmoiron/sqlx"
)

type UserRepository interface {
	GetByEmailOrUsername(email string, username string) (*entity.User, error)
	Create(user *entity.User) error
	GetByID(id int64) (*entity.User, error)
	List(query *model.PaginationQuery) ([]*entity.User, int64, error)
	Update(user *entity.User) error
	SoftDelete(id int64) error
}

type userRepository struct {
	db     *sqlx.DB
	logger *utils.Logger
}

func NewUserRepository(db *sqlx.DB, logger *utils.Logger) UserRepository {
	return &userRepository{db: db, logger: logger}
}

func (r *userRepository) GetByEmailOrUsername(email string, username string) (*entity.User, error) {
	query := `SELECT * FROM users WHERE usr_email = $1 OR usr_username = $2`
	user := &entity.User{}
	err := r.db.Get(user, query, email, username)
	if err != nil {
		r.logger.Error("UserRepository.GetByEmailOrUsername: %v", err)
		return nil, err
	}

	return user, nil
}

func (r *userRepository) Create(user *entity.User) error {
	query := `
		INSERT INTO users (usr_username, usr_password, usr_email, usr_created_at, usr_updated_at)
		VALUES ($1, $2, $3, NOW(), NOW()) RETURNING usr_id, usr_created_at, usr_updated_at
	`

	return r.db.QueryRowx(query, user.Username, user.Password, user.Email).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
}

func (r *userRepository) GetByID(id int64) (*entity.User, error) {
	user := &entity.User{}
	query := `SELECT * FROM users WHERE usr_id = $1 AND usr_deleted_at IS NULL`

	err := r.db.Get(user, query, id)
	if err != nil {
		r.logger.Error("UserRepository.GetByID: %v", err)
		return nil, err
	}

	return user, nil
}

func (r *userRepository) List(query *model.PaginationQuery) ([]*entity.User, int64, error) {
	var total int64
	countQuery := `SELECT COUNT(*) FROM users WHERE usr_deleted_at IS NULL`
	err := r.db.Get(&total, countQuery)
	if err != nil {
		r.logger.Error("UserRepository.List: %v", err)
		return nil, 0, err
	}

	users := []*entity.User{}
	listQuery := fmt.Sprintf(`
		SELECT * FROM users
		WHERE usr_deleted_at IS NULL
		ORDER BY usr_created_at DESC
		LIMIT %d OFFSET %d
	`, query.Limit, query.Offset)

	err = r.db.Select(&users, listQuery)
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *userRepository) Update(user *entity.User) error {
	query := `
		UPDATE users
		SET usr_username = $1, usr_email = $2, usr_updated_at = NOW()
		WHERE usr_id = $3
		RETURNING usr_updated_at
	`

	return r.db.QueryRowx(query, user.Username, user.Email, user.ID).Scan(&user.UpdatedAt)
}

func (r *userRepository) SoftDelete(id int64) error {
	query := `UPDATE users SET usr_deleted_at = NOW() WHERE usr_id = $1`
	_, err := r.db.Exec(query, id)
	if err != nil {
		r.logger.Error("UserRepository.Delete: %v", err)
	}
	return err
}
