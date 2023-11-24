package user

type Service interface {
	Create(request UserCreate) (error)
	Login(request UserLogin) (string, error)
	UpdatePassword(request UserUpdate) (error)
}
