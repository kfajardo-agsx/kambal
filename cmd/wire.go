//+build wireinject

package cmd

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/google/wire"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	// load postgres driver
	_ "github.com/lib/pq"
	"github.com/kfajardo-agsx/kambal.git/internal/component/account"

	// services
	dropboxSvc "github.com/kfajardo-agsx/kambal.git/internal/component/dropbox"
	fileSvc "github.com/kfajardo-agsx/kambal.git/internal/component/file"
	minioSvc "github.com/kfajardo-agsx/kambal.git/internal/component/minio"
	"github.com/kfajardo-agsx/kambal.git/internal/component/store"
	"github.com/kfajardo-agsx/kambal.git/internal/component/store_provider"
	tenantSvc "github.com/kfajardo-agsx/kambal.git/internal/component/tenant"

	"github.com/kfajardo-agsx/kambal.git/internal/infrastructure/dropbox"
	minioApp "github.com/kfajardo-agsx/kambal.git/internal/infrastructure/minio"
	"github.com/kfajardo-agsx/kambal.git/internal/infrastructure/postgres"
	"github.com/kfajardo-agsx/kambal.git/internal/infrastructure/postgres/repository"

	// controllers
	"github.com/kfajardo-agsx/kambal.git/internal/entrypoint/api/rest"
	accountHandler "github.com/kfajardo-agsx/kambal.git/internal/entrypoint/api/rest/account"
	dropboxHandler "github.com/kfajardo-agsx/kambal.git/internal/entrypoint/api/rest/dropbox"
	fileHandler "github.com/kfajardo-agsx/kambal.git/internal/entrypoint/api/rest/file"
	storeHandler "github.com/kfajardo-agsx/kambal.git/internal/entrypoint/api/rest/store"
	storeProviderHandler "github.com/kfajardo-agsx/kambal.git/internal/entrypoint/api/rest/store_provider"
	tenantHandler "github.com/kfajardo-agsx/kambal.git/internal/entrypoint/api/rest/tenant"
)

// this is where we wire all dependencies to run the API
type Keys struct {
	// APIKey      string `json:"api-key"`
	// S3AccessKey string `json:"s3-access-key"`
	// S3SecretKey string `json:"s3-secret-key"`

	DBUsername string `json:"db-username"`
	DBPassword string `json:"db-password"`
}

func createRestAPI() *rest.API {
	wire.Build(
		// Allow configuration of credentials from file
		ProvideKeysFromFile,
		ProvideRestAPIConfig,

		// database
		ProvideDatasource,
		ProvideGormDB,

		ProvideFileContext,

		// repository
		// repository.NewTenantRepository,
		// repository.NewStoreProviderRepository,
		// repository.NewAccountRepository,
		// repository.NewStoreRepository,
		// repository.NewFileRepository,
		repository.NewUserRepository,

		// ProvideDropboxConfig,
		// dropbox.NewService,
		// minioApp.NewMinioStorage,

		// service
		userSvc.NewUserService,

		// tenantSvc.NewTenantService,
		// store_provider.NewStoreProviderService,
		// account.NewAccountService,
		// store.NewStoreService,
		// fileSvc.NewFileService,
		// dropboxSvc.NewDropboxService,
		// minioSvc.NewMinioService,

		// controller
		userHandler.NewController,

		// tenantHandler.NewController,
		// storeProviderHandler.NewController,
		// accountHandler.NewController,
		// storeHandler.NewController,
		// fileHandler.NewController,
		// dropboxHandler.NewController,

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

func ProvideGormDB(datasource *postgres.Datasource) *gorm.DB {
	db, err := gorm.Open("postgres", datasource.AsPQString())
	if err != nil {
		log.WithError(err).Error("unable to get gorm db connection")
		os.Exit(1)
	}
	return db
}
