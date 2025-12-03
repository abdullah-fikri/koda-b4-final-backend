package models

import (
	"backend/config"
	"backend/utils"
	"context"
	"time"
)

type ListUserStruct struct {
	ID             int64  `json:"id"`
	Email          string `json:"email"`
	Role           string `json:"role"`
	Username       string `json:"username"`
	Phone          string `json:"phone"`
	Address        string `json:"address"`
	ProfilePicture string `json:"profile_picture"`
	CreatedAt time.Time `json:"since"`
}


type UpdateUserRequest struct {
    Username string `json:"username" binding:"required,min=3,max=20"`
    Phone    string `json:"phone" binding:"omitempty,min=10,max=15"`
    Address  string `json:"address" binding:"omitempty,max=100"`
    Password string `json:"password" binding:"omitempty,min=6,max=32"`
}


func GetUserProfile(userId int64)(ListUserStruct, error){
	ctx := context.Background()
	query := `
	SELECT 
	u.id,
	u.created_at,
	u.email,
	p.username,
	p.phone,
	p.address,
	COALESCE(p.profile_picture, '') AS profile_picture
	FROM users u
	LEFT JOIN profile p ON p.users_id = u.id
	WHERE u.id = $1`

	var u ListUserStruct
	err := config.Db.QueryRow(ctx, query, userId).Scan(
		&u.ID,
		&u.CreatedAt,
		&u.Email,
		&u.Username,
		&u.Phone,
		&u.Address,
		&u.ProfilePicture,
	)
	if err != nil{
		return ListUserStruct{}, err
	}
	return u, nil
}



func UpdateUserModel(id int64, req UpdateUserRequest) (*User, error) {
    ctx := context.Background()

    if req.Password != "" {
        hashedPassword := utils.HashPassword(req.Password)

        _, err := config.Db.Exec(ctx,
            `UPDATE users SET password=$1 WHERE id=$2`,
            hashedPassword, id,
        )
        if err != nil {
            return nil, err
        }
    }

    _, err := config.Db.Exec(ctx,
        `UPDATE profile 
         SET username=$1
         WHERE users_id=$2`,
        req.Username,id,
    )
    if err != nil {
        return nil, err
    }

    var user User
    err = config.Db.QueryRow(ctx,
        `SELECT id, email, role FROM users WHERE id=$1`,
        id,
    ).Scan(&user.ID, &user.Email, &user.Role)
    if err != nil {
        return nil, err
    }

    return &user, nil
}



func UpdateUserProfilePicture(userID int64, path string) error {
	ctx := context.Background()
	_, err := config.Db.Exec(ctx,
		`UPDATE profile 
		 SET profile_picture = $1, updated_at = NOW()
		 WHERE users_id = $2`,
		path, userID,
	)
	return err
}