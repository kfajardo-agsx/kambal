package common

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// WriteErrorResponse is a convenience function for writing error responses
func WriteErrorResponse(res http.ResponseWriter, status int, jsonErrs ...JSONError) {
	res.WriteHeader(status)
	if err := json.NewEncoder(res).Encode(jsonErrs); err != nil {
		log.WithError(err).Error("unable to send error response, the JSON error cound not be encoded")
	}
}

// WriteResponse is a convenience function for writing responses
func WriteResponse(res http.ResponseWriter, status int, entity interface{}) {
	res.WriteHeader(status)
	if err := json.NewEncoder(res).Encode(entity); err != nil {
		log.WithError(err).Error("unable to send response, the entity could not be encoded")
	}
}
