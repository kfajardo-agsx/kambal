package auth

import (
	"net/http"

	log "github.com/sirupsen/logrus"
	"gitlab.com/amihan/common/libraries/go/responses.git"
)

// A composable router middleware for using API keys to authorize callers.
// The presence and validty of the API key will override/skip all further security processing.
// The implication is that any bearer of the API key is trusted to make all calls.
// TODO: This probably belongs in a common library
type APIKeyAuthorize struct {
	// The value of the API key. For now it's a global key for the API.
	APIKey string
	// Wbere to look for the api key
	APIKeyParam string
}

type MiddlewareHandlerFunc func(next http.HandlerFunc) http.HandlerFunc

// Conditionally delegates to `standardAuthChain` if an API key is not found in the header.
// If an api key is found this will bypass the standard auth chain.
// If an api key is not found the element in standard auth chain will be chained into a call to resource.
// ex. AuthorizeOnAPIKeyMiddleware(a, b, c, d) will result in the following evaluation b(c(d(a)))
func (auth *APIKeyAuthorize) AuthorizeOnAPIKeyMiddleware(handler http.HandlerFunc, standardAuthChain ...MiddlewareHandlerFunc) http.HandlerFunc {
	// Recursively build standard auth chain for use when API key is not sent.
	// Note the arguments to the method are
	//	* resource - the handler that is to be put behind an auth chain.
	//   	If the API key set and valid the rest of the auth chain is ignored, the request is forwarded to this handler
	//	* chained - this is the built chain of handler functions. The chain will be executed in the same order as passed into standardAuthChain
	//  * authChain - the middleware functions to build up the chain.
	// 		This function calls itself recursively, each time taking the element at the end and attaching it to the beginning of the chain.
	//		This will result in the functions evaluated in the order they are here. eg [a, b, c] will be evaluated as a(b(c(resour
	var implem func(http.HandlerFunc, http.HandlerFunc, []MiddlewareHandlerFunc) http.HandlerFunc
	implem = func(resource http.HandlerFunc, chained http.HandlerFunc, authChain []MiddlewareHandlerFunc) http.HandlerFunc {
		chainLength := len(authChain)
		if chainLength == 0 {
			return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
				val := req.Header.Get(auth.APIKeyParam)
				if val == "" {
					log.Debugf("API Key not set in header %s", auth.APIKeyParam)
					chained(res, req)
					return
				}
				if val != auth.APIKey {
					log.Debugf("Provided api key does not match. Expected %s got %s", auth.APIKey, val)
					responses.WriteUnauthorizedError(res)
					return
				}
				resource(res, req)
			})
		}
		// Recursively walk the authChain slice backwards

		return implem(resource, authChain[chainLength-1](chained), authChain[0:chainLength-1])
	}
	return implem(handler, handler, standardAuthChain)
}

// Require the api key in the request.
func (auth *APIKeyAuthorize) AuthorizeOnRequiredApiKeyMiddleware(handler http.HandlerFunc) http.HandlerFunc {
	return auth.AuthorizeOnAPIKeyMiddleware(handler,
		func(h http.HandlerFunc) http.HandlerFunc {
			return func(res http.ResponseWriter, req *http.Request) {
				log.Debugf("The API Key is required. Request not authorized.")
				responses.WriteUnauthorizedError(res)
			}
		})
}
