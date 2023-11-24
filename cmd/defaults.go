package cmd

import (
	"github.com/spf13/viper"
)

var defaults = map[string]interface{}{
	"log.level": "debug",

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
	// Secrets file for setting credentials in environments
	"secrets.file":       "config/secrets.json",
	"file.max-upload-mb": 25,

	// database configuration
	"datasource.type":       "postgres",
	"datasource.host":       "localhost",
	"datasource.port":       5432,
	"datasource.database":   "files",
	"datasource.sslMode":    "disable",
	"datasource.migrations": "db/migrations",

	// Basic JWT auth config
	"api.rest.auth.jwtPubKeyUrl":     "https://<PUT-KEYCLOAK-HOSTNAME-HERE>/auth/realms/<PUT-REALM-NAME-HERE>",
	"api.rest.auth.claimsAttribute":  "<PUT-ATTRIBUTE-IN-CLAIM-THAT-IDENTIFIES-THE-USER-HERE>",
	"api.rest.auth.requestParamName": "<PUT-NAME-OF-PARAMETER-THAT-IDENTIFIES-THE-OWNER-OF-RESOURCES-HERE>",

	// Auth via a passed API key
	"api.rest.auth.apiKeyParamName": "API-KEY",

	// RBAC configuration
	"api.rest.auth.rbacFile": "config/rbac.yaml",

	"dropbox.redirectUri": "http://localhost:8080/api/v1/authorize/redirect",
	"dropbox.auth.url":    "https://www.dropbox.com",
	"dropbox.api.url":     "https://api.dropboxapi.com",
	"dropbox.content.url": "https://content.dropboxapi.com",
}

func init() {
	for key, value := range defaults {
		viper.SetDefault(key, value)
	}
}
