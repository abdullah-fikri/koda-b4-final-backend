package models

import (
	"backend/config"
	"backend/utils"
	"context"
)

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
	Username string `json:"username" binding:"required"`
	Phone    string `json:"phone" binding:"omitempty"`
	Address  string `json:"address" binding:"omitempty"`
	Role     string `json:"-"`
}


type User struct {
	ID       int64  `json:"id"`
	Email    string `json:"email"`
	Password string `json:"-"`
	Role     string `json:"role"`
}


type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}


func Register(req RegisterRequest) (*User, error) {
	ctx := context.Background()

	hashedPassword := utils.HashPassword(req.Password)

	var userID int64
	err := config.Db.QueryRow(ctx,
		`INSERT INTO users (email, password, role)
         VALUES ($1, $2, 'user')
         RETURNING id`,
		req.Email, hashedPassword,
	).Scan(&userID)
	if err != nil {
		return nil, err
	}

	_, err = config.Db.Exec(ctx,
		`INSERT INTO profile (users_id, username, phone, address)
         VALUES ($1, $2, $3, $4)`,
		userID, req.Username, req.Phone, req.Address,
	)
	if err != nil {
		return nil, err
	}

	return &User{
		ID:    userID,
		Email: req.Email,
		Role:  "user",
	}, nil
}


// login 
func Login(email string) (*User, error) {
	ctx := context.Background()

	var user User
	err := config.Db.QueryRow(ctx,
		`SELECT id, email, password, role
		 FROM users
		 WHERE email = $1`,
		email,
	).Scan(&user.ID, &user.Email, &user.Password, &user.Role)

	if err != nil {
		return nil, err
	}

	return &user, nil
}



func Forgot(email string) (*User, error) {
	ctx := context.Background()

	var user User
	err := config.Db.QueryRow(ctx,
		`SELECT id,email,role FROM users WHERE email=$1`, email).Scan(&user.ID, &user.Email, &user.Role)

	if err != nil {
		return nil, err
	}

	return &user, nil
}



func UpdateUserPassword(email, hashedPassword string) error {
	ctx := context.Background()

	_, err := config.Db.Exec(ctx,
		`UPDATE users SET password=$1 WHERE email=$2`,
		hashedPassword, email,
	)

	return err
}
