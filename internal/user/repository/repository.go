package repository

import (
	"context"
	"database/sql"
	"go_project_template/configs/db"
	"go_project_template/internal/user/model"
)

type IUserRepository interface {
	GetUsers(ctx context.Context) ([]model.User, error)
	GetUserById(ctx context.Context, id int64) (model.User, error)
	AddNewUser(ctx context.Context, newUser model.User) (int64, error)
	UpdateUser(ctx context.Context, user model.User) (int64, error)
	DeleteUser(ctx context.Context, id int64) error
	UpdateEmailVerificationStatus(ctx context.Context, email string) error
	FindByEmail(ctx context.Context, email string) (model.User, error)
}

type UserRepository struct {
	db db.DBInterface
}

func NewUserRepository(db db.DBInterface) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

// Query Implementation

func (q *UserRepository) GetUsers(ctx context.Context) ([]model.User, error) {
	var users []model.User

	sqlStatement := `
	SELECT * FROM public.users
	`

	// Querying
	rows, err := q.db.QueryContext(ctx, sqlStatement)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var user model.User

		// read data
		err := rows.Scan(&user.ID, &user.Fullname)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

func (q *UserRepository) GetUserById(ctx context.Context, id int64) (model.User, error) {
	var user model.User

	queryStatement := `
	SELECT
		user_id,
		fullname,
		email,
		created_at,
		updated_at
	FROM "user".users
	WHERE user_id=$1
	`
	err := q.db.QueryRowContext(ctx, queryStatement, id).Scan(
		&user.ID,
		&user.Fullname,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return user, err
	}

	return user, nil
}

func (q *UserRepository) AddNewUser(ctx context.Context, newUser model.User) (int64, error) {

	var newId int64

	sqlStatement := `
	INSERT INTO
    "user".users(fullname, email, password)
	VALUES
			($1, $2, $3)
	RETURNING user_id
	`

	err := q.db.QueryRowContext(ctx, sqlStatement, newUser.Fullname, newUser.Email, newUser.Password).Scan(&newId)

	if err != nil {
		return 0, err
	}

	return newId, nil
}

func (q *UserRepository) UpdateUser(ctx context.Context, user model.User) (int64, error) {

	sqlStatement := `
	UPDATE public.users SET name=$2 WHERE id=$1
	`

	res, err := q.db.ExecContext(ctx, sqlStatement, user.ID, user.Fullname)

	if err != nil {
		return 0, err
	}

	rowsAffeced, err := res.RowsAffected()

	if err != nil {
		return 0, err
	}
	return rowsAffeced, nil
}

func (q *UserRepository) DeleteUser(ctx context.Context, id int64) error {

	// create sql statement to delete user from database
	sqlStatement := `DELETE FROM public.users WHERE id=$1`

	// execute sql statement
	res, err := q.db.ExecContext(ctx, sqlStatement, id)

	if err != nil {
		return err
	}

	rowsAffected, _ := res.RowsAffected()

	if rowsAffected < 1 {
		return sql.ErrNoRows
	}
	return nil
}

func (q *UserRepository) UpdateEmailVerificationStatus(ctx context.Context, email string) error {

	updateStatement := `
	UPDATE
		"user".users
	SET 
		is_verified = TRUE
	WHERE
		users.email = $1
	`

	res, err := q.db.ExecContext(ctx, updateStatement, email)

	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil

}

func (q *UserRepository) FindByEmail(ctx context.Context, email string) (model.User, error) {
	var user model.User

	queryStatement := `
		SELECT 
			user_id,
			fullname,
			email,
			is_verified,
			created_at,
			updated_at
		FROM "user".users
		WHERE
			email=$1
	`

	err := q.db.QueryRowContext(ctx, queryStatement, email).Scan(
		&user.ID,
		&user.Fullname,
		&user.Email,
		&user.IsVerified,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return user, err
	}

	return user, nil
}
