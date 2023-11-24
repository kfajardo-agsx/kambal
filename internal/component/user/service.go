package user

import (
	log "github.com/sirupsen/logrus"
	"github.com/kfajardo-agsx/kambal.git/internal/component/common"
)

type UserService struct {
	repository    Repository
	config UserConfig
}

type UserConfig struct {
	SecretKey string
	DefaultPassword string
}

// NewUserService
func NewUserService(repo Repository) Service {
	return &UserService{
		repository:    repo,
	}
}

func (s *UserService) Create(request UserCreate) (error) {
	if (request.Password == "") {
		request.Password = s.config.DefaultPassword
	}

	encrypted := common.Encrypt(s.config.SecretKey, request.Password)
	userDB := User{
		Username: request.Username,
		EncryptedPassword: encrypted,
		UserRole: request.UserRole,
		CustomerID: request.CustomerID,
	}
	created, err := s.repository.Create(userDB)
	if err != nil {
		return common.RepositoryErrorToAPIError(err)
	}
	log.Info("========================")
	log.Info("USER CREATED")
	log.Info(created)
	log.Info("========================")
	return nil
}
func (s *UserService) Login(request UserLogin) (string, error) {
	userAcct, err := s.repository.Get(request.Username)
	if err != nil  {
		return "", err
	}
	encryptedPassword := common.Encrypt(s.config.SecretKey, request.Password)
	if encryptedPassword != userAcct.EncryptedPassword {
		return "", common.NewAPIError(common.ErrorTypeUnauthorized, "Username/Password incorrect")
	}
	return userAcct.ID, nil
}

func (s *UserService) UpdatePassword(request UserUpdate) (error) {
	userAcct, err := s.repository.Get(request.Username)
	if err != nil  {
		return err
	}
	encryptedPassword := common.Encrypt(s.config.SecretKey, request.OldPassword)
	if encryptedPassword != userAcct.EncryptedPassword {
		return common.NewAPIError(common.ErrorTypeUnauthorized, "Old Password incorrect")
	}
	err = s.repository.UpdatePassword(request.Username, request.NewPassword)
	if err != nil {
		return common.RepositoryErrorToAPIError(err)
	}
	return nil
}