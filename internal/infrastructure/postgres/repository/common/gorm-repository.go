package common

import (
	"github.com/jinzhu/gorm"
)

type GormRepository struct {
	*gorm.DB
}
