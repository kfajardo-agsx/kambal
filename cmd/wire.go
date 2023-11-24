//go:build wireinject
// +build wireinject

package cmd

import (
	"encoding/json"
	"io/ioutil"
	"os"

	// load postgres driver
	_ "github.com/lib/pq"
	"gitlab.com/amihan/core/base.git/internal/component/user"
	"gitlab.com/amihan/core/base.git/internal/entrypoint/api/rest"
	"gitlab.com/amihan/core/base.git/internal/infrastructure/postgres"
	"gitlab.com/amihan/core/base.git/internal/infrastructure/postgres/repository"

	handler "gitlab.com/amihan/core/base.git/internal/entrypoint/api/rest/base"

	"github.com/google/wire"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// this is where we wire all dependencies to run the API
type Keys struct {
	DBUsername string `json:"db-username"`
	DBPassword string `json:"db-password"`
	APIKey     string `json:"api-key"`
}

func createRestAPI() *rest.API {
	wire.Build(
		ProvideKeysFromFile,
		ProvideRestAPIConfig,

		ProvideDatasource,
		ProvideGormDB,
		repository.NewGormUserRepository,

		// services
		user.NewUserService,

		handler.NewController,

		mux.NewRouter,
		rest.NewRestAPI,
	)
	return &rest.API{}
}

func createMigration() *postgres.Migration {
	wire.Build(
		ProvideKeysFromFile,
		ProvideDatasource,
		postgres.NewMigration,
	)
	return &postgres.Migration{}
}

func ProvideKeysFromFile() *Keys {
	file := viper.GetString("secrets.file")
	jsonFile, err := os.Open(file)
	if err != nil {
		log.WithError(err).Error("unable to read security keys file")
		os.Exit(1)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var secrets Keys
	json.Unmarshal(byteValue, &secrets)

	return &secrets
}

func ProvideRestAPIConfig(keys *Keys) *rest.Config {
	var config rest.Config
	err := viper.UnmarshalKey("api.rest", &config)
	if err != nil {
		log.WithError(err).Error("unable to read RestAPIConfig")
		os.Exit(1)
	}

	rbacData, err := ioutil.ReadFile(config.Auth.RBACFile)
	if err != nil {
		log.WithError(err).Error("Error in reading RBAC data")
	}

	config.Auth.RBAC = string(rbacData)
	config.Auth.APIKey = keys.APIKey

	config.Version = root.Version

	log.Info("========================================")
	log.Info("API Configuration")
	log.Info("========================================")
	log.Info("Host:    ", config.Host)
	log.Info("Port:    ", config.Port)
	log.Info("Spec:    ", config.Spec)
	log.Info("Version: ", config.Port)

	log.Debugf("API Config: %+v", config)
	level := viper.GetBool("log.debug")
	if level {
		log.SetLevel(logrus.DebugLevel)
	}
	return &config
}

func ProvideDatasource(keys *Keys) *postgres.Datasource {
	var datasource postgres.Datasource
	err := viper.UnmarshalKey("datasource", &datasource)
	if err != nil {
		log.WithError(err).Error("unable to read Datasource config")
		os.Exit(1)
	}
	datasource.Username = keys.DBUsername
	datasource.Password = keys.DBPassword
	return &datasource
}

func ProvideGormDB(datasource *postgres.Datasource) *gorm.DB {
	db, err := gorm.Open("postgres", datasource.AsPQString())
	if err != nil {
		log.WithError(err).Error("unable to get gorm db connection")
		os.Exit(1)
	}
	return db
}
