package base

import (
	"net/http"

	"github.com/gorilla/mux"
	middleware "gitlab.com/amihan/common/libraries/go/jwt-auth.git"
	"gitlab.com/amihan/core/base.git/internal/component/user"
	"gitlab.com/amihan/core/base.git/internal/infrastructure/auth"
)

type Controller struct {
	userService user.Service
}

func NewController(userService user.Service) *Controller {
	return &Controller{
		userService: userService,
	}
}

func (c *Controller) Register(router *mux.Router, apiKeyAuthMiddleware *auth.APIKeyAuthorize, jwtMiddleware *middleware.JWT) {
	base := router.PathPrefix("/users").Subrouter()
	userService := c.userService

	base.
		Methods(http.MethodPost, http.MethodOptions).
		Path("").
		HandlerFunc(apiKeyAuthMiddleware.AuthorizeOnRequiredApiKeyMiddleware(
			createUser(userService)))

}
