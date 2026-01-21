package repository

import (
	"Lolopenza/CRUD-F/internal/models"
	"context"
	"database/sql"
)

type Repo struct {
	DB *sql.DB
}

func New(db *sql.DB) *Repo {
	return &Repo{DB: db}
}

func (r *Repo) CreateUser(ctx context.Context, email, name, surname string) (int, error) {
	var id int

	stmt := `INSERT INTO users ( email, name, surname) VALUES ($1, $2, $3) RETURNING usr_id`

	err := r.DB.QueryRowContext(ctx, stmt, email, name, surname).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, err
}

func (r *Repo) GetAllUsers(ctx context.Context) ([]models.User, error) {
	rows, err := r.DB.QueryContext(ctx, `SELECT usr_id, email, name, surname FROM users`)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var users []models.User

	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.User_ID, &u.Email, &u.Name, &u.Surname); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}

func (r *Repo) GetUserByID(ctx context.Context, num_id int) (models.User, error) {
	var u models.User

	query := `SELECT usr_id, email, name, surname
			FROM users
			WHERE usr_id = $1`

	err := r.DB.QueryRowContext(ctx, query, num_id).Scan(&u.User_ID, &u.Email, &u.Name, &u.Surname)
	if err != nil {
		return u, err
	}

	return u, nil
}

func (r *Repo) UpdateUser(ctx context.Context, num_id int, email, name, surname string) (models.User, error) {
	var u models.User

	query := `
		UPDATE users
		SET email = $1, name = $2, surname = $3, updated_at = now()
		WHERE usr_id = $4
		RETURNING usr_id, email, name, surname, created_at, updated_at
	`

	err := r.DB.QueryRowContext(ctx, query, email, name, surname, num_id).
		Scan(&u.User_ID, &u.Email, &u.Name, &u.Surname, &u.Created_At, &u.Updated_At)

	if err != nil {
		return u, err
	}

	return u, nil
}

func (r *Repo) DeleteUser(ctx context.Context, num_id int) error {
	query := `
			DELETE FROM users
			WHERE usr_id = $1
			`

	res, err := r.DB.ExecContext(ctx, query, num_id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return err
}
