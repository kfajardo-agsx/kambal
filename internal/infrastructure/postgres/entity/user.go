package entity

type User struct {
	Base
	Username  string `gorm:"username"`
	EncryptedPassword string `gorm:"encrypted_password"`
	UserRole    string `gorm:"user_role"`
	CustomerID string `gorm:"customer_id"`
}

func (repoModel *User) GetID() interface{} {
	return repoModel.ID
}