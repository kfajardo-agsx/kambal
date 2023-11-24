package middlewares

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"gitlab.com/amihan/common/libraries/go/responses.git"

	"github.com/gorilla/mux"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
	log "github.com/sirupsen/logrus"
)

type (
	// openIDConfiguration a subset of the openid-configuration that has the issuer and jwks_url only.
	openIDConfiguration struct {
		Issuer  string `json:"issuer"`
		JWKsURI string `json:"jwks_uri"`
	}

	// openIDMiddleware is the actual middleware that handles the JWT verification
	openIDMiddleware struct {
		jwkSet *jwk.Set
		issuer string
		next   http.Handler
	}
)

// OpenID is an OpenID middleware that adds JWT validation
func OpenID(wellKnownConfiguration string) (middleware mux.MiddlewareFunc, err error) {
	// first we get the openid-configuration
	config, err := getOpenIDConfiguration(wellKnownConfiguration)
	if err != nil {
		log.WithError(err).Error("unable to get wellknownconfiguration")
		return
	}
	// then we get the JWKs
	jwkSet, err := getJWK(config.JWKsURI)
	if err != nil {
		log.WithError(err).Error("unable to get JWK")
		return
	}
	middleware = func(next http.Handler) http.Handler {
		return &openIDMiddleware{
			jwkSet: jwkSet,
			issuer: config.Issuer,
			next:   next,
		}
	}
	return
}

// ServeHTTP runs the middleware filter to verify the JWT signature, issuer, and validity
func (m *openIDMiddleware) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	token := extractToken(req)
	if token == "" {
		responses.WriteUnauthorizedError(res)
		return
	}
	if err := verify(token, m.issuer, m.jwkSet.Keys); err != nil {
		responses.WriteUnauthorizedError(res)
		return
	}
	m.next.ServeHTTP(res, req)
}

// retrieve OpenID Configuration from the provided URL.
func getOpenIDConfiguration(wellKnownConfiguration string) (config *openIDConfiguration, err error) {
	resp, err := http.Get(wellKnownConfiguration)
	if err != nil {
		log.WithError(err).Error("unable to get well known configuration from the provided URL")
		return
	}
	body := resp.Body
	defer body.Close()
	rawConfig, err := ioutil.ReadAll(body)
	if err != nil {
		log.WithError(err).Error("unable to read the response from the provided URL")
		return
	}
	if err = json.Unmarshal(rawConfig, &config); err != nil {
		log.WithError(err).Error("unable to read the response from the provided URL as JSON")
		return
	}
	log.WithFields(log.Fields{
		"issuer":   config.Issuer,
		"jwks_url": config.JWKsURI,
	}).Info("OpenID Configuration retrieved")
	return
}

// extract the token from the auth header. Format should be = Authorization: Bearer {token}. should the token be not
// in the proper format or the extract token is not JWT, this will return an empty string.
func extractToken(req *http.Request) string {
	authHeader := req.Header.Get("Authorization")
	bearer := strings.Split(authHeader, " ")
	if len(bearer) != 2 || bearer[0] != "Bearer" {
		return ""
	}
	rawToken := bearer[1]
	tokenSegments := len(strings.Split(rawToken, "."))
	if tokenSegments != 3 {
		return ""
	}
	return rawToken
}

// get the JWKs information for verifying signatures
func getJWK(jwksURL string) (*jwk.Set, error) {
	set, err := jwk.Fetch(jwksURL)
	if err != nil {
		log.WithError(err).Error("unable to fetch JWK")
		return nil, err
	}
	return set, nil
}

// verify the token. this checks the signature, the issuer, and expiry.
func verify(rawToken, issuer string, keys []jwk.Key) error {

	var token *jwt.Token
	var err error
	for _, key := range keys {
		s := jwa.SignatureAlgorithm(key.Algorithm())
		k, e := key.Materialize()
		if e != nil {
			log.WithError(e).Warning("unable to materialze key, skipping")
			continue
		}
		token, err = jwt.ParseString(rawToken, jwt.WithVerify(s, k))
		if err == nil {
			break
		}
	}
	if token == nil || err != nil {
		return err
	}

	// // manual verification
	sameIssuer := issuer == token.Issuer()
	if !sameIssuer {
		return errors.New("Token validation error. Issuer is invalid")
	}
	expired := token.Expiration().Before(time.Now())
	if expired {
		return errors.New("Token validation error. Expired")
	}
	return nil
}
