package cmd

import (
	"github.com/spf13/viper"
)

var defaults = map[string]interface{}{
	"api.rest.host":                "0.0.0.0",
	"api.rest.port":                8080,
	"api.rest.spec":                "./openapi.yaml",
	"api.rest.cors.allowedOrigins": []string{"*"},
	"api.rest.cors.allowedHeaders": []string{
		"Content-Type",
		"Sec-Fetch-Dest",
		"Referer",
		"accept",
		"Sec-Fetch-Mode",
		"Sec-Fetch-Site",
		"User-Agent",
		"API-KEY",
		"Authorization",
	},
	"api.rest.cors.allowedMethods": []string{
		"OPTIONS",
		"GET",
		"POST",
		"PUT",
		"DELETE",
	},

	"api.rest.auth.jwtPubKeyUrl": "http://keycloak.dev.msme.amihan.net/auth/realms/bdo-msme-customer-app",
	"api.rest.auth.rbacFile":     "config/rbac.yaml",

	"api.rest.auth.apiKeyParamName": "API-KEY",

	"log.debug": true,

	"datasource.type":       "postgres",
	"datasource.host":       "localhost",
	"datasource.port":       5432,
	"datasource.database":   "loans",
	"datasource.sslMode":    "disable",
	"datasource.migrations": "db/migrations",

	"secrets.file": "config/secrets.json",
}

func init() {
	for key, value := range defaults {
		viper.SetDefault(key, value)
	}
}
