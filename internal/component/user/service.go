package user

type UserService struct {
	repository Repository
}

// NewUserService creates the default implementation of the user service
func NewUserService(repository Repository) Service {
	return &UserService{
		repository: repository,
	}
}
