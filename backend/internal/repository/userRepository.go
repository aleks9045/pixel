package repository

import (
	"database/sql"
	"errors"

	"github.com/aleks9045/pixels/internal/model"
)

type UserRepository struct {
	conn *sql.DB
}

func NewUserRepository(conn *sql.DB) *UserRepository {
	return &UserRepository{conn: conn}
}

func (r *UserRepository) ExistsByUserId(userId string) (bool, error) {
	var exists bool
	err := r.conn.QueryRow(
		`SELECT EXISTS(SELECT 1 FROM users WHERE user_id=$1)`, userId).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *UserRepository) GetIsAdmin(userId string) (bool, error) {
	var isAdmin bool
	err := r.conn.QueryRow(
		`SELECT is_admin FROM users WHERE user_id=$1`, userId).Scan(&isAdmin)
	if err != nil {
		return false, err
	}
	return isAdmin, nil
}

func (r *UserRepository) UpdateUser(userModel *model.User) error {
	_, err := r.conn.Exec(
		`UPDATE users
		SET name = $1, student_group = $2, time_limit = $3
		WHERE user_id = $4`,
		userModel.Name,
		userModel.StudentGroup,
		userModel.TimeLimit,
		userModel.UserId,
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) GetUsers(page, size int) ([]*model.User, error) {
	offset := page * size

	rows, err := r.conn.Query(`SELECT user_id, username, name, avatar_url, student_group, time_limit, is_admin
		 FROM users
		 OFFSET $1
		 LIMIT $2`, offset, size)
	defer rows.Close()
	if err != nil {
		return []*model.User{}, err
	}
	var (
		userId       string
		userName     string
		name         string
		avatarURL    string
		studentGroup sql.NullString
		timeLimit    int32
		isAdmin      bool
	)
	users := make([]*model.User, size)
	i := 0
	for rows.Next() {
		if err := rows.Scan(&userId, &userName, &name, &avatarURL, &studentGroup, &timeLimit, &isAdmin); err != nil {
			return []*model.User{}, err
		}
		user := &model.User{
			UserId:       userId,
			Username:     userName,
			Name:         name,
			AvatarUrl:    avatarURL,
			StudentGroup: studentGroup.String,
			TimeLimit:    timeLimit,
			IsAdmin:      isAdmin,
		}
		users[i] = user
		i++
	}
	if rows.Err() != nil {
		return []*model.User{}, err
	}
	cutIdx := 0
	for _, u := range users {
		if u == nil {
			break
		}
		cutIdx++
	}

	return users[:cutIdx], nil

}

func (r *UserRepository) CreateIfNotExists(user *model.OauthUser, userId string) (*string, error) {
	var existingUUID string

	err := r.conn.QueryRow(
		"SELECT user_id FROM users WHERE username = $1",
		user.Username,
	).Scan(&existingUUID)

	if err == nil {
		return &existingUUID, nil
	}

	if !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	_, err = r.conn.Exec(
		"INSERT INTO users (user_id, username, name, avatar_url) VALUES ($1, $2, $3, $4)",
		userId,
		user.Username,
		user.Name,
		user.AvatarURL,
	)
	if err != nil {
		return nil, err
	}

	return &userId, nil
}

func (r *UserRepository) GetByUserId(userId string) (*model.User, error) {
	var (
		userName     string
		name         string
		timeLimit    int32
		avatarURL    string
		studentGroup sql.NullString
		isAdmin      bool
	)

	err := r.conn.QueryRow(
		`SELECT username, name, avatar_url, student_group, time_limit, is_admin
		 FROM users
		 WHERE user_id = $1`,
		userId,
	).Scan(
		&userName,
		&name,
		&avatarURL,
		&studentGroup,
		&timeLimit,
		&isAdmin,
	)

	if err != nil {
		return nil, err
	}

	return &model.User{
		UserId:       userId,
		Username:     userName,
		Name:         name,
		AvatarUrl:    avatarURL,
		StudentGroup: studentGroup.String,
		TimeLimit:    timeLimit,
		IsAdmin:      isAdmin,
	}, nil
}
