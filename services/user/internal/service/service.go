package service

import (
	"context"
	"regexp"
	"shopik/services/user/internal/domain"
	"shopik/services/user/internal/password"
)

var (
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	hasLower   = regexp.MustCompile(`[a-z]`)
	hasUpper   = regexp.MustCompile(`[A-Z]`)
	hasDigit   = regexp.MustCompile(`[0-9]`)
	hasSpecial = regexp.MustCompile(`[!@#$%^&*().,?":{}|<>]`)
)

type UserRepository interface {
	Save(ctx context.Context, user domain.User) (id int64, err error)
	GetByEmail(ctx context.Context, email string) (user domain.User, err error)
}

type TokenIssuer interface {
	Issuer(userID string) (string, error)
}

type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func IsValidEmail(email string) bool {
	return emailRegex.MatchString(email)
}
func IsValidPassword(password string) bool {
	if len(password) < 8 || len(password) > 128 {
		return false
	}

	if !hasLower.MatchString(password) {
		return false
	}

	if !hasUpper.MatchString(password) {
		return false
	}

	if !hasDigit.MatchString(password) {
		return false
	}

	if !hasSpecial.MatchString(password) {
		return false
	}

	return true
}

func (s *UserService) RegisterNewUser(ctx context.Context, registerInput domain.RegisterInput) (id int64, err error) {
	if !IsValidEmail(registerInput.Email) {
		return 0, domain.ErrInvalidEmail
	}
	if !IsValidPassword(registerInput.Password) {
		return 0, domain.ErrInvalidPassword
	}
	hashPass, err := password.HashPassword(registerInput.Password)
	if err != nil {
		return 0, err
	}

	user := domain.User{
		Email:        registerInput.Email,
		PasswordHash: hashPass,
	}

	id, err = s.repo.Save(ctx, user)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *UserService) Login(ctx context.Context, loginInput domain.LoginInput) (jwt string, err error) {
	if !IsValidEmail(loginInput.Email) {
		return "", domain.ErrInvalidEmail
	}
	if !IsValidPassword(loginInput.Password) {
		return "", domain.ErrInvalidPassword
	}
	user, err := s.repo.GetByEmail(ctx, loginInput.Email)
	if err != nil {
		return "", domain.ErrInvalidCredentials
	}
	ok, err := password.ComparePassword(loginInput.Password, user.PasswordHash)
	if err != nil {
		return "", err
	}
	if !ok {
		return "", domain.ErrInvalidCredentials
	}

}
