package model

import (
	"github.com/jmoiron/sqlx"
	"time"
)

type User struct {
	Id        uint64    `db:"id"`
	Username  string    `db:"username"`
	Password  string    `db:"password"`
	Salt      string    `db:"salt"`
	Lang      string    `db:"lang"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
type UserRepo struct {
	conn *sqlx.DB
}

func NewUserRepo(conn *sqlx.DB) *UserRepo {
	return &UserRepo{
		conn: conn,
	}
}

func (t *UserRepo) Get(id int) (User, error) {
	var user User
	err := t.conn.Get(&user, "select * from user where id = ?", id)
	return user, err
}

func (t *UserRepo) GetByUsername(username string) (User, error) {
	var user User
	err := t.conn.Get(&user, "select * from user where username = ?", username)
	return user, err
}
