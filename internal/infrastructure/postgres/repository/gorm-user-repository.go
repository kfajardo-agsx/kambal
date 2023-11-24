package repository

import (
	// "fmt"

	// "github.com/jinzhu/copier"
	"github.com/jinzhu/gorm"
	"github.com/kfajardo-agsx/kambal.git/internal/component/user"
	"github.com/kfajardo-agsx/kambal.git/internal/infrastructure/postgres/entity"
	"github.com/kfajardo-agsx/kambal.git/internal/infrastructure/postgres/repository/common"
)

type UserRepository struct {
	*common.GormRepository
}

func NewUserRepository(db *gorm.DB) user.Repository {
	return &UserRepository{
		&common.GormRepository{
			db.Debug(),
		},
	}
}


func (repo *UserRepository) Get(username string) (*user.User, error) {
	data := &user.User{}
	err := repo.FindEntity(&entity.User{}, data, "username = ?", username)
	if nil != err {
		return nil, err
	}

	return data, nil
}
	
func (repo *UserRepository) UpdatePassword(username string, newPassword string) (error) {
	data, err := repo.Get(username)
	if err != nil {
		return err
	}
	data.EncryptedPassword = newPassword

	_, err = repo.Create(*data)
	return err
}

func (repo *UserRepository) Create(body user.User) (*user.User, error) {
	_, err := repo.SaveEntity(&entity.User{}, body)
	if nil != err {
		return nil, err
	}

	return repo.Get(body.Username)
}

func (repo *UserRepository) Delete(username string) error {
	return repo.DeleteEntity(entity.User{}, "username = ?", username)
}
