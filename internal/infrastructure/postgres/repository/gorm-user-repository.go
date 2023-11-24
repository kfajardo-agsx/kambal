package repository

import (
	// "strconv"
	// "time"

	"github.com/jinzhu/gorm"
	"gitlab.com/amihan/core/base.git/internal/component/user"
	"gitlab.com/amihan/core/base.git/internal/infrastructure/postgres/repository/common"
)

type GormUserRepository struct {
	*common.GormRepository
}

func NewGormUserRepository(db *gorm.DB) user.Repository {
	return &GormUserRepository{
		&common.GormRepository{
			db.Debug(),
		},
	}
}
