package services

import (
	"auth_service/internal/models"
	"auth_service/pkg/jwt"
	"context"
	"database/sql"
	"errors"
	"time"


	// "github.com/redis/go-redis/v9"
	// //"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo  UserRepository
	redis RedisClient
	email EmailClient
}
type UserRepository interface {
	CreateUser(user *models.User) error
	GetUserByEmail(email string) (*models.User, error)
	GetUserByID(id string) (*models.User, error)
}

type RedisClient interface {
	SetSession(ctx context.Context, sessionID, userID string, ttl time.Duration) error
	GetSession(ctx context.Context, sessionID string) (string, error)
	DeleteSession(ctx context.Context, sessionID string) error
}
type EmailClient interface {
	SendEmail(ctx context.Context, to, subject, body string) error
}

func NewAuthService(repo UserRepository, redis RedisClient, email EmailClient) *AuthService {
	return &AuthService{
		repo:  repo,
		redis: redis,
		email: email,
	}
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
		return nil, "", errors.New("invalid email or password")
	}

	//session id
	token, err := jwt.GenerateToken(user.ID)
	if err != nil {
		return nil, "", err
	}

	//Store in redis
	err = s.redis.SetSession(ctx, token, user.ID, 24*time.Hour)
	if err != nil {
		return nil, "Failed to create session", err
	}

	return user, token, nil
}

//Get me

func (s *AuthService) GetMe(ctx context.Context, sessionId string) (*models.User, error) {
	//get user_id from redis
	userID := sessionId 
	
	//fetch from db
	user, err := s.repo.GetUserByID(userID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *AuthService) Logout(ctx context.Context, sessionID string) error {
	return s.redis.DeleteSession(ctx, sessionID)
}



func (s *AuthService) SendEmail(ctx context.Context, sessionID string) error {

	// 1. Validate session
	userID := sessionID 
	
	// 2. Get user
	user, err := s.repo.GetUserByID(userID)
	if err != nil {
		return err
	}

	// 3. Call email service
	return s.email.SendEmail(
		ctx,
		user.Email,
		"Welcome!",
		"Hello " + user.FirstName + ", this is your email 🎉",
	)
}
