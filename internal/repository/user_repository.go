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
	tx, err := r.db.Beginx()
	if err != nil {
		r.logger.Error("UserRepository.GetByEmailOrUsername: failed to start transaction: %v", err)
		return nil, err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	query := `SELECT * FROM users WHERE usr_email = $1 OR usr_username = $2 LIMIT 1`
	user := &entity.User{}

	err = tx.Get(user, query, email, username)
	if err != nil {
		r.logger.Error("UserRepository.GetByEmailOrUsername: %v", err)
		return nil, err
	}

	return user, nil
}

func (r *userRepository) Create(user *entity.User) error {
	tx, err := r.db.Beginx()
	if err != nil {
		r.logger.Error("UserRepository.Create: failed to start transaction: %v", err)
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	query := `
		INSERT INTO users (usr_username, usr_password, usr_email, usr_created_at, usr_updated_at)
		VALUES ($1, $2, $3, NOW(), NOW()) RETURNING usr_id, usr_created_at, usr_updated_at
	`

	err = tx.QueryRowx(query, user.Username, user.Password, user.Email).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		r.logger.Error("UserRepository.Create: %v", err)
		return err
	}

	if err := tx.Commit(); err != nil {
		r.logger.Error("UserRepository.Create: failed to commit transaction: %v", err)
		return err
	}

	return nil
}

func (r *userRepository) GetByID(id int64) (*entity.User, error) {
	tx, err := r.db.Beginx()
	if err != nil {
		r.logger.Error("UserRepository.GetByID: failed to start transaction: %v", err)
		return nil, err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	user := &entity.User{}
	query := `SELECT * FROM users WHERE usr_id = $1 AND usr_deleted_at IS NULL`

	err = tx.Get(user, query, id)
	if err != nil {
		r.logger.Error("UserRepository.GetByID: %v", err)
		return nil, err
	}

	return user, nil
}

func (r *userRepository) List(query *model.PaginationQuery) ([]*entity.User, int64, error) {
	tx, err := r.db.Beginx()
	if err != nil {
		r.logger.Error("UserRepository.List: failed to start transaction: %v", err)
		return nil, 0, err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	var total int64
	countQuery := `SELECT COUNT(*) FROM users WHERE usr_deleted_at IS NULL`
	err = tx.Get(&total, countQuery)
	if err != nil {
		r.logger.Error("UserRepository.List: %v", err)
		return nil, 0, err
	}
	r.logger.Info("UserRepository.List: executed query: %v", countQuery)

	users := []*entity.User{}
	listQuery := fmt.Sprintf(`
		SELECT * FROM users
		WHERE usr_deleted_at IS NULL
		ORDER BY usr_created_at DESC
		LIMIT %d OFFSET %d
	`, query.Limit, query.Offset)

	err = tx.Select(&users, listQuery)
	if err != nil {
		return nil, 0, err
	}
	r.logger.Info("UserRepository.List: executed query: %v", listQuery)

	return users, total, nil
}

func (r *userRepository) Update(user *entity.User) error {
	tx, err := r.db.Beginx()
	if err != nil {
		r.logger.Error("UserRepository.Update: failed to start transaction: %v", err)
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	query := `
		UPDATE users
		SET usr_username = $1, usr_email = $2, usr_updated_at = NOW()
		WHERE usr_id = $3
		RETURNING usr_updated_at
	`

	err = tx.QueryRowx(query, user.Username, user.Email, user.ID).Scan(&user.UpdatedAt)
	if err != nil {
		r.logger.Error("UserRepository.Update: %v", err)
		return err
	}

	if err := tx.Commit(); err != nil {
		r.logger.Error("UserRepository.Update: failed to commit transaction: %v", err)
		return err
	}

	return nil
}

func (r *userRepository) SoftDelete(id int64) error {
	tx, err := r.db.Beginx()
	if err != nil {
		r.logger.Error("UserRepository.SoftDelete: failed to start transaction: %v", err)
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	query := `UPDATE users SET usr_deleted_at = NOW() WHERE usr_id = $1`
	_, err = tx.Exec(query, id)
	if err != nil {
		r.logger.Error("UserRepository.SoftDelete: %v", err)
		return err
	}
	r.logger.Info("UserRepository.SoftDelete: executed query: %v", query)

	if err := tx.Commit(); err != nil {
		r.logger.Error("UserRepository.SoftDelete: failed to commit transaction: %v", err)
		return err
	}

	return err
}
