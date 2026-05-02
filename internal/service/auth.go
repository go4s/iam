package service

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/go4s/iam/internal/repository"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = func() []byte {
	if secret := os.Getenv("JWT_SECRET"); secret != "" {
		return []byte(secret)
	}
	panic("should never reach here, JWT_SECRET is required in environment")
}()

type AuthService struct {
	userRepo *repository.UserRepository
}

func NewAuthService(userRepo *repository.UserRepository) *AuthService {
	return &AuthService{userRepo: userRepo}
}

func (s *AuthService) Login(username, password string) (string, map[string]any, error) {
	user, err := s.userRepo.GetByUsername(username)
	if err != nil {
		return "", nil, err
	}
	if user == nil {
		return "", nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", nil, errors.New("invalid credentials")
	}

	if user.Status == "disabled" {
		return "", nil, errors.New("user account is disabled")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":     user.Username,
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", nil, err
	}

	// 构建用户基本信息
	userInfo := map[string]any{
		"id":       fmt.Sprintf("user:%d", user.ID),
		"username": user.Username,
	}

	return tokenString, userInfo, nil
}

func (s *AuthService) Me(username string) (map[string]any, error) {
	user, err := s.userRepo.GetByUsername(username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	// 获取用户角色
	userRoleRepo := &repository.UserRoleRepository{}
	roleIDs, err := userRoleRepo.GetRoleIDsByUserID(user.ID)
	if err != nil {
		return nil, err
	}

	roleRepo := &repository.RoleRepository{}
	var roles []string
	for _, rid := range roleIDs {
		role, err := roleRepo.GetByID(rid)
		if err != nil || role == nil {
			continue
		}
		roles = append(roles, fmt.Sprintf("role:%s", role.Code))
	}

	return map[string]any{
		"id":         fmt.Sprintf("user:%d", user.ID),
		"username":   user.Username,
		"roles:":     roles,
		"roles-fold": false,
	}, nil
}

func (s *AuthService) ParseToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid claims")
	}

	return claims, nil
}

func GetJWTSecret() []byte {
	return jwtSecret
}
