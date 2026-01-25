package response

import (
	"time"
	"wallet_api/internal/entity"
)

type UserResponse struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	CreatedAt string `json:"created_at"`
}

func ToUserDto(user *entity.User) UserResponse {
	return UserResponse{
		ID:        user.ID.String(),
		Username:  user.Username,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
	}
}
