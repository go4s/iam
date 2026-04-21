package service

import (
    "errors"
    "os"
    "time"
    
    "github.com/go4s/iam/internal/auth"
    "github.com/go4s/iam/internal/model"
    "github.com/go4s/iam/internal/repository"
    
    "github.com/golang-jwt/jwt/v5"
    "golang.org/x/crypto/bcrypt"
)

var jwtSecret = func() []byte {
    if secret := os.Getenv("JWT_SECRET"); secret != "" {
        return []byte(secret)
    }
    return []byte("your-secret-key") // Default for development
}()

type AuthService struct {
    userRepo *repository.UserRepository
}

func NewAuthService(userRepo *repository.UserRepository) *AuthService {
    return &AuthService{userRepo: userRepo}
}

func (s *AuthService) Register(username, password, role string) error {
    existing, err := s.userRepo.GetByUsername(username)
    if err != nil {
        return err
    }
    if existing != nil {
        return errors.New("user already exists")
    }
    
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return err
    }
    
    user := &model.User{
        Username:     username,
        PasswordHash: string(hashedPassword),
        Role:         role,
    }
    
    if err := s.userRepo.Create(user); err != nil {
        return err
    }
    
    // Add role to Casbin grouping policy
    _, err = auth.Enforcer.AddGroupingPolicy(username, role)
    return err
}

func (s *AuthService) Login(username, password string) (string, error) {
    user, err := s.userRepo.GetByUsername(username)
    if err != nil {
        return "", err
    }
    if user == nil {
        return "", errors.New("invalid credentials")
    }
    
    if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
        return "", errors.New("invalid credentials")
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "sub":  user.Username,
        "role": user.Role,
        "exp":  time.Now().Add(time.Hour * 24).Unix(),
    })
    
    return token.SignedString(jwtSecret)
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
