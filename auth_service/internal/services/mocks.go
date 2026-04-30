package services

import (
	"auth_service/internal/models"
	"context"
	"database/sql"
	"time"
)

type MockRepo struct {
	users map[string]*models.User
}

func (m *MockRepo) CreateUser(user *models.User) error {
	m.users[user.Email] = user
	return nil
}

func (m *MockRepo) GetUserByEmail(email string) (*models.User, error) {
	if user, ok := m.users[email]; ok {
		return user, nil
	}
	return nil, sql.ErrNoRows
}

func (m *MockRepo) GetUserByID(id string) (*models.User, error) {
	for _, u := range m.users {
		if u.ID == id {
			return u, nil
		}
	}
	return nil, sql.ErrNoRows
}

type MockRedis struct {
	store map[string]string
}

func (m *MockRedis) SetSession(ctx context.Context, sessionID, userID string, ttl time.Duration) error {
	m.store[sessionID] = userID
	return nil
}

func (m *MockRedis) GetSession(ctx context.Context, sessionID string) (string, error) {
	if userID, ok := m.store[sessionID]; ok {
		return userID, nil
	}
	return "", sql.ErrNoRows
}

func (m *MockRedis) DeleteSession(ctx context.Context, sessionID string) error {
	delete(m.store, sessionID)
	return nil
}