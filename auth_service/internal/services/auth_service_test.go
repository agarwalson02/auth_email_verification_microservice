package services

import (
	"auth_service/internal/models"
	"context"
	"testing"
)

func TestRegister(t *testing.T) {
	repo := &MockRepo{users: make(map[string]*models.User)}
	redis := &MockRedis{store: make(map[string]string)}

	svc := NewAuthService(repo, redis)

	user := &models.User{
		Email:     "test@example.com",
		FirstName: "Test",
		LastName:  "User",
		Password:  "123456",
	}

	created, err := svc.Register(user)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if created.Email != user.Email {
		t.Fatalf("expected email %s, got %s", user.Email, created.Email)
	}
}

func TestRegisterDuplicate(t *testing.T) {
	repo := &MockRepo{users: make(map[string]*models.User)}
	redis := &MockRedis{store: make(map[string]string)}

	svc := NewAuthService(repo, redis)

	user := &models.User{
		Email: "test@example.com",
	}

	_, _ = svc.Register(user)

	_, err := svc.Register(user)

	if err == nil {
		t.Fatal("expected error for duplicate email")
	}
}

func TestLoginSuccess(t *testing.T) {
	repo := &MockRepo{users: make(map[string]*models.User)}
	redis := &MockRedis{store: make(map[string]string)}

	svc := NewAuthService(repo, redis)

	// register first
	user := &models.User{
		Email:    "test@example.com",
		Password: "123456",
	}
	_, _ = svc.Register(user)

	ctx := context.Background()

	u, session, err := svc.Login(ctx, "test@example.com", "123456")

	if err != nil {
		t.Fatalf("expected success, got error %v", err)
	}

	if session == "" {
		t.Fatal("expected session id")
	}

	if u.Email != "test@example.com" {
		t.Fatal("wrong user returned")
	}
}

func TestLoginFail(t *testing.T) {
	repo := &MockRepo{users: make(map[string]*models.User)}
	redis := &MockRedis{store: make(map[string]string)}

	svc := NewAuthService(repo, redis)

	ctx := context.Background()

	_, _, err := svc.Login(ctx, "no@example.com", "123")

	if err == nil {
		t.Fatal("expected error")
	}
}