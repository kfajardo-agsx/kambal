package user

type Repository interface {
	Get(username string) (*User, error)
	UpdatePassword(username string, newPassword string) (error)
	Create(user User) (*User, error)
}
