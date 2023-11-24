package user

import (
	"net/http"

	"github.com/gorilla/mux" // router/dispatcher
	// apiKeyAuth "gitlab.com/amihan/common/libraries/go/api-key-auth.git"
	// middleware "gitlab.com/amihan/common/libraries/go/jwt-auth.git" // claims
	"github.com/kfajardo-agsx/kambal.git/internal/component/user"
)

type Controller struct {
	userService user.Service
}

func NewController(userService user.Service) *Controller {
	return &Controller{
		userService: userService,
	}
}

func (c *Controller) Register(router *mux.Router) {
	application := router.PathPrefix("/users").Subrouter()

	application.
		Methods(http.MethodPost, http.MethodOptions).
		Path("").
		HandlerFunc(create(c.userService))

	application.
		Methods(http.MethodPost, http.MethodOptions).
		Path("/login").
		HandlerFunc(login(c.userService))

	application.
		Methods(http.MethodPost, http.MethodOptions).
		Path("/update").
		HandlerFunc(update(c.userService))

}
