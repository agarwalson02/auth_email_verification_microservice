package services

import (
	"auth_service/internal/models"
	"auth_service/internal/repository"
	"auth_service/pkg/redis"
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"

	// "github.com/redis/go-redis/v9"
	// //"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo  *repository.UserRepository
	redis *redis.Client
}

func NewAuthService(repo *repository.UserRepository, redis *redis.Client) *AuthService {
	return &AuthService{repo: repo, redis: redis}
}

func (s *AuthService) Register(user *models.User) (*models.User, error) {

	//check user exsist
	exsistUser, _ := s.repo.GetUserByEmail(user.Email)
	if exsistUser != nil {
		return nil, errors.New("Email already exists")
	}

	//hashing password
	hashedpassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		return nil, err
	}

	user.Password = string(hashedpassword)

	//save user
	err = s.repo.CreateUser(user)
	if err != nil {
		return nil, err
	}
	return user, nil

}

// login
func (s *AuthService) Login(ctx context.Context, email, password string) (*models.User, string, error) {
	// find user
	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, "", errors.New("invalid email or password")
		}
		return nil, "", err
	}
	// password compare
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, "Invalid email or password", nil
	}

	//session id
	sessionId := uuid.New().String()

	//Store in redis
	err = s.redis.SetSession(ctx, sessionId, user.ID, 24*time.Hour)
	if err != nil {
		return nil, "Failed to create session", err
	}

	return user, sessionId, nil
}

//Get me

func (s *AuthService) GetMe(ctx context.Context, sessionId string) (*models.User, error) {
	//get user_id from redis
	userId, err := s.redis.GetSession(ctx, sessionId)
	if err != nil {
		return nil, errors.New("invalid session")
	}

	//fetch from db
	user, err := s.repo.GetUserById(userId)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *AuthService) Logout(ctx context.Context, sessionID string) error {
	return s.redis.DeleteSession(ctx, sessionID)
}
