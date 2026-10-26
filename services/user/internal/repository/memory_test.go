package repository_test

import (
	"fmt"
	"shopik/services/user/internal/domain"
	"shopik/services/user/internal/repository"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSaveSuccess(t *testing.T) {
	repo := repository.NewInMemoryRepository()
	wantID := int64(1)
	user := domain.User{
		Email:        "balbes@gmail.com",
		PasswordHash: "$argon2id$v=19$m=65536,t=1,p=4$6gGS5CiokEYOQrmfqHki1A$mLlUyoa3xBcUU3WiPEYnnH8l2Gi7C4qYYQugJcXuPiA",
		CreatedAt:    time.Now(),
	}
	id, err := repo.Save(t.Context(), user)
	assert.NoError(t, err)
	assert.Equal(t, wantID, id)
}

func TestSaveSpamSuccess(t *testing.T) {
	var wg sync.WaitGroup
	repo := repository.NewInMemoryRepository()
	for i := 0; i < 40; i++ {
		wg.Add(1)
		go func() {
			user := domain.User{
				Email:        fmt.Sprintf("balbes%d@gmail.com", i),
				PasswordHash: "$argon2id$v=19$m=65536,t=1,p=4$6gGS5CiokEYOQrmfqHki1A$mLlUyoa3xBcUU3WiPEYnnH8l2Gi7C4qYYQugJcXuPiA",
				CreatedAt:    time.Now(),
			}
			_, err := repo.Save(t.Context(), user)
			assert.NoError(t, err)
			wg.Done()
		}()
	}
	wg.Wait()
	for i := 0; i < 40; i++ {
		email := fmt.Sprintf("balbes%d@gmail.com", i)
		user, err := repo.GetByEmail(t.Context(), email)
		assert.NoError(t, err)
		assert.Equal(t, user.Email, email)
	}
}

func TestGetByEmailSuccess(t *testing.T) {
	repo := repository.NewInMemoryRepository()
	email := "balbes@gmail.com"
	passwordHash := "$argon2id$v=19$m=65536,t=1,p=4$6gGS5CiokEYOQrmfqHki1A$mLlUyoa3xBcUU3WiPEYnnH8l2Gi7C4qYYQugJcXuPiA"
	createdAt := time.Now()
	user := domain.User{
		Email:        email,
		PasswordHash: passwordHash,
		CreatedAt:    createdAt,
	}
	wantUser := user
	wantUser.ID = 1
	repo.Save(t.Context(), user)
	gotUser, err := repo.GetByEmail(t.Context(), email)
	assert.NoError(t, err)
	assert.Equal(t, wantUser, gotUser)
}
