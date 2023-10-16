package repository

import (
	"context"
)

func (r *Repository) GetTestById(ctx context.Context, input GetTestByIdInput) (output GetTestByIdOutput, err error) {
	err = r.Db.QueryRowContext(ctx, "SELECT name FROM test WHERE id = $1", input.Id).Scan(&output.Name)
	if err != nil {
		return
	}
	return
}

func (r *Repository) CreateUser(ctx context.Context, user User) (id int, err error) {
	err = r.Db.QueryRowContext(ctx, "INSERT INTO users (full_name, phone_number, password, created_at, updated_at, login_count) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id", user.FullName, user.PhoneNumber, user.Password, user.CreatedAt, user.UpdatedAt, user.LoginCount).Scan(&id)
	if err != nil {
		return
	}
	return
}

func (r *Repository) GetUserByPhoneNumber(ctx context.Context, phoneNumber string) (user User, err error) {
	err = r.Db.QueryRowContext(ctx, "SELECT id, full_name, phone_number, password, created_at, updated_at, last_login, login_count FROM users WHERE phone_number = $1", phoneNumber).Scan(&user.Id, &user.FullName, &user.PhoneNumber, &user.Password, &user.CreatedAt, &user.UpdatedAt, &user.LastLogin, &user.LoginCount)
	if err != nil {
		return
	}
	return
}

func (r *Repository) UpdateUser(ctx context.Context, user User) (err error) {
	_, err = r.Db.ExecContext(ctx, "UPDATE users SET full_name = $1, phone_number = $2, password = $3, created_at = $4, updated_at = $5, last_login = $6, login_count = $7 WHERE id = $8", user.FullName, user.PhoneNumber, user.Password, user.CreatedAt, user.UpdatedAt, user.LastLogin, user.LoginCount, user.Id)
	if err != nil {
		return
	}
	return
}

func (r *Repository) GetUserById(ctx context.Context, id int) (user User, err error) {
	err = r.Db.QueryRowContext(ctx, "SELECT id, full_name, phone_number, password, created_at, updated_at, last_login, login_count FROM users WHERE id = $1", id).Scan(&user.Id, &user.FullName, &user.PhoneNumber, &user.Password, &user.CreatedAt, &user.UpdatedAt, &user.LastLogin, &user.LoginCount)
	if err != nil {
		return
	}
	return
}
