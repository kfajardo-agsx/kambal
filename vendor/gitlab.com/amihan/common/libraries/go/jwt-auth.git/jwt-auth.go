package jwt

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
	"gitlab.com/amihan/common/libraries/go/responses.git"
)

type Exception struct {
	Message string `json:"message"`
}

type RBACItem struct {
	Role      string     `yaml:"role"`
	Endpoints []Endpoint `yaml:"endpoints"`
}

type Endpoint struct {
	Method string `yaml:"method"`
	Path   string `json:"path"`
}

type JWT struct {
	RBAC      []RBACItem `yaml:"rbac"`
	PubKeyURL string
}

type UserGroup struct {
	Groups []string `json:"groups"`
}

func WriteJWTError(res http.ResponseWriter, err error) {
	status := http.StatusBadRequest
	newErr := responses.ErrorBuilder().
		Status(status).
		Title("JWT Validation Error").
		Detail(fmt.Sprintf("%v", err)).
		Build()
	responses.WriteErrorResponse(res, status, newErr)
}

func (j *JWT) ValidateMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		tokenHeader := req.Header.Get("Authorization")
		log.Debugf("Token Header: %+v", tokenHeader)

		if tokenHeader == "" {
			log.Error("JWT token is missing or invalid")
			WriteJWTError(w, errors.New("JWT Token is missing or invalid"))
			return
		}

		if strings.Index(tokenHeader, "Bearer") != 0 {
			log.Error("JWT token is missing or invalid")
			WriteJWTError(w, errors.New("JWT Token is missing or invalid"))
			return
		}

		accessToken := strings.Split(tokenHeader, " ")[1]
		splitted := strings.Split(accessToken, ".")
		if len(splitted) != 3 {
			log.Error("JWT token is missing or invalid")
			WriteJWTError(w, errors.New("JWT Token is missing or invalid"))
			return
		}

		keyFunc := func(token *jwt.Token) (interface{}, error) {
			log.Debug("getting key for token: ", token)
			log.Debug("raw: ", token.Raw)
			log.Debug("method: ", token.Method)
			log.Debug("header: ", token.Header)
			log.Debug("claims: ", token.Claims)
			log.Debug("signature: ", token.Signature)
			log.Debug("valid: ", token.Valid)

			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return false, errors.New("Error in signing JWT")
			}

			//Getting public key from keycloak server
			client := http.Client{}
			url := j.PubKeyURL

			res, err := client.Get(url)

			if err != nil {
				return nil, err
			}

			defer res.Body.Close()

			respString, err := ioutil.ReadAll(res.Body)

			if err != nil {
				return nil, err
			}

			// Convert object to interface
			var respObj map[string]interface{}
			err = json.Unmarshal([]byte(respString), &respObj)

			pubKey := "-----BEGIN PUBLIC KEY-----\n" + respObj["public_key"].(string) + "\n-----END PUBLIC KEY-----"
			verifyKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(pubKey))

			return verifyKey, nil
		}

		token, err := jwt.Parse(accessToken, keyFunc)
		if err != nil {
			log.Error("Token parsing failed: ", err)
			WriteJWTError(w, errors.New("Token parsing failed"))
			return
		}

		if !token.Valid {
			log.Error("Invalid Token")
			WriteJWTError(w, errors.New("Invalid Token"))
			return
		}

		url := req.URL
		method := req.Method

		if err := j.jwtAuthorize(token, url.String(), method); err != nil {
			log.Error("User is not authorized to access this endpoint")
			WriteJWTError(w, errors.New("User is not authorized to access this endpoint"))
			return
		}

		ctx := req.Context()
		ctx = context.WithValue(ctx, "claims", token.Claims)
		r := req.WithContext(ctx)
		next(w, r)
	})
}

func (j *JWT) jwtAuthorize(token *jwt.Token, url string, method string) error {
	var ug UserGroup
	mapstructure.Decode(token.Claims, &ug)

	success := j.isAuthorized(ug.Groups, url, method)
	if !success {
		log.Error("User is not authorize to access this endpoint")
		return errors.New("User is not authorize to access this endpoint")
	}
	log.Debugf("JWTAuthorize is OK")
	return nil
}

func (j *JWT) isAuthorized(groups []string, url string, method string) bool {
	log.Debugf("RBAC VALUES %+v", j.RBAC)
	log.Debugf("URL Value %+v", url)
	log.Debugf("Method Value %+v", method)
	flag := false
	for index, element := range j.RBAC {
		// searching for role
		log.Debugf("INDEX %d ELEMENT %+v", index, element)
		for _, user := range groups {
			log.Debugf("user value is %+v", user)
			if user == element.Role {
				log.Debugf("user.Role == %s", element.Role)
				for _, el := range element.Endpoints {
					log.Debugf("EL VALUE IS %s", el)
					if el.Method == method {
						log.Debugf("el.Method value is %s", el.Method)
						if el.Path == url {
							log.Debugf("el.Path value is %s", el.Path)
							return true
						} else if el.Path != url {
							matched, err := regexp.MatchString(el.Path, url)
							if err != nil {
								log.Error("Error while matching regex value: ", err)
							}
							if matched {
								return true
							}
						} else {
							flag = false
							log.Error("Unauthorized")
						}
					}
				}

			}
		}
	}
	return flag
}
