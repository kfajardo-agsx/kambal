package user

type User struct {
	ID string `json:"id"`
	Username string `json:"username"`
	EncryptedPassword string `json:"password"`
	UserRole string `json:"user_role"`
	CustomerID string `json:"customer_id"`
}

type UserCreate struct {
	Username string `json:"username"`
	Password string `json:"password"`
	UserRole string `json:"user_role"`
	CustomerID string `json:"customer_id"`
}

type UserLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserUpdate struct {
	Username string `json:"username"`
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}
