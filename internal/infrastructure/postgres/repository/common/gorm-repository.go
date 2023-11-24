package common

import (
	"time"

	"github.com/jinzhu/copier" // copier library
	"github.com/jinzhu/gorm"   // gorm library for database
	log "github.com/sirupsen/logrus"
	"github.com/kfajardo-agsx/kambal.git/internal/infrastructure/postgres/entity"
)

type GormRepository struct {
	*gorm.DB
}

func (cr *GormRepository) SaveEntity(e entity.Entity, data interface{}) (id interface{}, err error) {
	if err = copier.Copy(e, data); err != nil {
		log.WithError(err).Error("unable to copy data to entity")
		return
	}

	err = cr.Save(e).Error
	id = e.GetID()
	return
}

func (cr *GormRepository) FindEntity(e interface{}, o interface{}, where ...interface{}) error {
	if err := cr.Find(e, where...).Error; err != nil {
		log.WithError(err).Error("error accessing database")
		return err
	}
	if err := copier.Copy(o, e); err != nil {
		log.WithError(err).Error("unable to copy entity to data")
	}
	return nil
}

func (cr *GormRepository) DeleteEntity(e interface{}, where ...interface{}) error {
	return cr.Delete(e, where...).Error
}

// func (cr *GormRepository) GetAccountID(externalID string) (string, error) {
// 	var accountEntity entity.Account
// 	err := cr.Table("accounts").Where("external_reference = ?", externalID).Find(&accountEntity).Error
// 	if err != nil {
// 		return "", err
// 	}
// 	return accountEntity.ID, nil
// }

func AddLimit(chain *gorm.DB, limit int) *gorm.DB {
	if limit == 0 {
		return chain
	}
	chain = chain.Limit(limit)
	return chain
}

func AddOffset(chain *gorm.DB, limit int, page int) *gorm.DB {
	if limit == 0 {
		return chain
	}

	if page == 0 {
		page = 1
	}

	chain = chain.Offset(limit * (page - 1))
	return chain
}

func AddDateCreatedFilter(chain *gorm.DB, startDate string, endDate string) *gorm.DB {
	dateLayout := "2006-01-02"
	t := time.Now().UTC()
	start, err := time.Parse(dateLayout, startDate)
	if err != nil {
		start = time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.UTC)
	}
	end, err := time.Parse(dateLayout, endDate)
	if err != nil {
		end = start.AddDate(0, 1, 0).Add(time.Nanosecond * -1)
	} else {
		end = end.AddDate(0, 0, 1).Add(time.Nanosecond * -1)
	}

	chain = chain.Where("created_at >= ?", start).Where("created_at <= ?", end)
	return chain
}

func AddOrdering(chain *gorm.DB, parameter string, ascending bool) *gorm.DB {
	if parameter == "" {
		parameter = "created_at"
	}
	if ascending {
		chain = chain.Order(parameter + " asc")
	} else {
		chain = chain.Order(parameter + " desc")
	}

	return chain
}

func AddLikeFilter(chain *gorm.DB, parameter string, value string) *gorm.DB {
	chain = chain.Where(parameter+" ILIKE ?", "%"+value+"%")
	return chain
}

func AddEqualFilter(chain *gorm.DB, parameter string, value string) *gorm.DB {
	chain = chain.Where(parameter+" = ?", value)
	return chain
}

// func (cr *GormRepository) GetTenantID(externalReference string) (string, error) {
// 	var tenantEntity entity.Tenant
// 	err := cr.Table("tenants").Where("external_reference = ?", externalReference).Find(&tenantEntity).Error
// 	if err != nil {
// 		return "", err
// 	}
// 	return tenantEntity.ID, nil
// }
