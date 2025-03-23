package entity

import (
	"time"
)

type User struct {
	ID        int64     `db:"usr_id"`
	Username  string    `db:"usr_username"`
	Email     string    `db:"usr_email"`
	Password  string    `db:"usr_password"`
	CreatedAt time.Time `db:"usr_created_at"`
	UpdatedAt time.Time `db:"usr_updated_at"`
	DeletedAt *time.Time `db:"usr_deleted_at"`
}

func (u *User) TableName() string {
	return "users"
}
