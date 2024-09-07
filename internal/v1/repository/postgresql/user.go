package postgresql

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"mini-network/internal/v1/models"
)

type UserRepo struct {
	connect *pgx.Conn
}

func NewUserRepo(connect *pgx.Conn) *UserRepo {
	return &UserRepo{connect}
}

func (r *UserRepo) GetUserByID(ctx context.Context, id int) (models.User, error) {
	var user models.User
	rows, err := r.connect.Query(ctx, `SELECT * FROM users WHERE ID =$1 LIMIT 1`, id)
	if err == nil {
		defer rows.Close()
		user, err = pgx.CollectOneRow(rows, pgx.RowToStructByName[models.User])
		if errors.Is(err, pgx.ErrNoRows) {
			err = nil
		}
	}
	return user, err
}

func (r *UserRepo) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	var user models.User
	rows, err := r.connect.Query(ctx, `SELECT * FROM users WHERE email =$1 LIMIT 1`, email)
	if err == nil {
		defer rows.Close()
		user, err = pgx.CollectOneRow(rows, pgx.RowToStructByName[models.User])
		if errors.Is(err, pgx.ErrNoRows) {
			err = nil
		}
	}
	return user, err
}

func (r *UserRepo) SaveUser(ctx context.Context, email, name, password string) (int, error) {
	var userId int
	err := r.connect.QueryRow(ctx, `INSERT INTO users(email,name,password) VALUES($1,$2,$3) ON CONFLICT("email") DO NOTHING RETURNING id`, email, name, password).Scan(&userId)
	return userId, err
}
