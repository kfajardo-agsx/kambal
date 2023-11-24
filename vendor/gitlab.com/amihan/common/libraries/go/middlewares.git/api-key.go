package middlewares

import (
	"net/http"

	"github.com/gorilla/mux"

	"gitlab.com/amihan/common/libraries/go/responses.git"
)

// APIKeyType is a custom type defining the location in a request where the API key can be found.
type APIKeyType string

const (
	// APIKeyTypeHeader means the API Key is found in the request header
	APIKeyTypeHeader APIKeyType = "HEADER"
	// APIKEYTypeQuery means the API Key is found in the query parameter
	APIKEYTypeQuery APIKeyType = "QUERY"
)

type apiKeyMiddleware struct {
	t APIKeyType
	k string
	v string
	h http.Handler
}

// APIKey returns a middleware that handles API Key validation.
//
// Parameters are:
// - keyType:        defines where the API Key is located in the request
// - parameterName:  the name of the parameter holding the key
// - parameterValue: the actual API Key
//
// If the API Key validation fails, this will return HTTP 401.
func APIKey(keyType APIKeyType, parameterName, parameterValue string) mux.MiddlewareFunc {
	apiKeyHandler := &apiKeyMiddleware{
		t: keyType,
		k: parameterName,
		v: parameterValue,
	}
	return func(handler http.Handler) http.Handler {
		apiKeyHandler.h = handler
		return apiKeyHandler
	}
}

// ServeHTTP contains the core logic of the middleware.
func (m *apiKeyMiddleware) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	value := ""
	switch m.t {
	case APIKeyTypeHeader:
		value = req.Header.Get(m.k)
	case APIKEYTypeQuery:
		value = req.URL.Query().Get(m.k)
	default:
		responses.WriteUnauthorizedError(res)
		return
	}
	if value != m.v {
		responses.WriteUnauthorizedError(res)
		return
	}
	m.h.ServeHTTP(res, req)
}
