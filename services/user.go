package services

import (
	"fmt"

	"github.com/example/supabase-migration-demo/models"
	"github.com/example/supabase-migration-demo/repositories"
	"github.com/google/uuid"
)

type UserService struct {
	repo *repositories.UserRepository
}

func NewUserService(repo *repositories.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser(req *models.CreateUserRequest) (*models.User, error) {
	existing, err := s.repo.GetByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, fmt.Errorf("user with email %s already exists", req.Email)
	}

	user := &models.User{
		Name:  req.Name,
		Email: req.Email,
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) GetAllUsers() ([]models.User, error) {
	return s.repo.GetAll()
}

func (s *UserService) GetUserByID(id uuid.UUID) (*models.User, error) {
	return s.repo.GetByID(id)
}

func (s *UserService) DeleteUser(id uuid.UUID) error {
	return s.repo.Delete(id)
}