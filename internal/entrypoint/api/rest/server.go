package rest

import (
	"fmt"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	middleware "gitlab.com/amihan/common/libraries/go/jwt-auth.git"
	"gitlab.com/amihan/common/libraries/go/middlewares.git"
	"gitlab.com/amihan/common/libraries/go/responses.git"
	application "gitlab.com/amihan/core/base.git/internal/entrypoint/api/rest/base"
	"gitlab.com/amihan/core/base.git/internal/infrastructure/auth"
	"gopkg.in/yaml.v2"
)

type (
	Config struct {
		Host    string
		Port    int
		Spec    string
		Version string
		Cors    CORSConfig
		Auth    AuthConfig
	}
	AuthConfig struct {
		APIKey          string
		APIKeyParamName string
		RBACFile        string
		RBAC            string
		JWTPubKeyURL    string
	}
	CORSConfig struct {
		AllowedOrigins []string
		AllowedHeaders []string
		AllowedMethods []string
	}

	API struct {
		config                *Config
		router                *mux.Router
		applicationController *application.Controller
	}
)

func NewRestAPI(config *Config, router *mux.Router, applicationController *application.Controller) *API {
	return &API{
		config:                config,
		router:                router,
		applicationController: applicationController,
	}
}

func (api *API) Run() error {
	api.router = api.router.PathPrefix("/api/v1").Subrouter()
	api.registerHandlers()
	api.exposeSwagger()
	api.exposeVersion()
	api.exposeHealth()
	api.addMiddlewares()
	api.enableCORS()
	return http.ListenAndServe(api.address(), api.router)
}

func (api *API) address() string {
	return fmt.Sprintf("%s:%d", api.config.Host, api.config.Port)
}

func (api *API) exposeSwagger() {
	api.router.HandleFunc("/spec", func(res http.ResponseWriter, req *http.Request) {
		http.ServeFile(res, req, api.config.Spec)
	})
	log.Infof("OpenAPI Spec accessible at http://%s/api/v1/spec", api.address())
}

func (api *API) exposeVersion() {
	api.router.HandleFunc("/version", func(res http.ResponseWriter, req *http.Request) {
		responses.WriteOKWithEntity(res, api.config.Version)
	})
}

func (api *API) exposeHealth() {
	api.router.HandleFunc("/health", func(res http.ResponseWriter, req *http.Request) {
		responses.WriteOK(res)
	})
}

func (api *API) enableCORS() {
	cors := handlers.CORS(
		handlers.AllowedOrigins(api.config.Cors.AllowedOrigins),
		handlers.AllowedHeaders(api.config.Cors.AllowedHeaders),
		handlers.AllowedMethods(api.config.Cors.AllowedMethods),
	)

	api.router.Use(cors)
	log.Info("CORS filter enabled")
}

func (api *API) addMiddlewares() {
	logger := middlewares.Logger(log.StandardLogger())
	api.router.Use(logger)
	log.Info("Logger filter enabled")
}

func (api *API) registerHandlers() {
	log.Infof("Register Handlers")
	jwtMiddleware := &middleware.JWT{
		PubKeyURL: api.config.Auth.JWTPubKeyURL,
	}

	apiKeyAuthMiddleware := &auth.APIKeyAuthorize{
		APIKey:      api.config.Auth.APIKey,
		APIKeyParam: api.config.Auth.APIKeyParamName,
	}

	err := yaml.Unmarshal([]byte(api.config.Auth.RBAC), &jwtMiddleware.RBAC)
	if err != nil {
		log.Errorf("Error decoding RBAC: %v", err.Error())
	}

	api.applicationController.Register(api.router, apiKeyAuthMiddleware, jwtMiddleware)
}
