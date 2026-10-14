package repository

import (
	"context"
	"sync"
)

type UserData struct {
	ID       int
	Username string
	PassHash string
}

type InMemoryRepository struct {
	userData map[int]UserData
	mu       sync.Mutex
	nextID   int
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		userData: make(map[int]UserData, 10),
		nextID:   1,
	}
}

func (r *InMemoryRepository) Save(ctx context.Context, user UserData) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	id := r.nextID
	r.userData[id] = user
	r.nextID++
	return id, nil
}
