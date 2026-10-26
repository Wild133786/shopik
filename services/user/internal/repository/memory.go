package repository

import (
	"context"
	"shopik/services/user/internal/domain"
	"sync"
)

type InMemoryRepository struct {
	userData map[string]domain.User
	mu       sync.Mutex
	nextID   int64
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		userData: make(map[string]domain.User, 10),
		nextID:   1,
	}
}

func (r *InMemoryRepository) Save(ctx context.Context, user domain.User) (int64, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	_, ok := r.userData[user.Email]
	if ok {
		return 0, domain.ErrEmailTaken
	}
	user.ID = r.nextID
	r.userData[user.Email] = user
	r.nextID++
	return user.ID, nil
}

func (r *InMemoryRepository) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	user, ok := r.userData[email]
	if !ok {
		return domain.User{}, domain.ErrUserNotFound
	}
	return user, nil
}
