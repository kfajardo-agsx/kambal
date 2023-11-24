package user

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/kfajardo-agsx/kambal.git/internal/component/user"
	"github.com/kfajardo-agsx/kambal.git/internal/entrypoint/api/rest/common"
)


func create(userSvc user.Service) http.HandlerFunc {
	log.Debug("handler registered: accounts::create")

	return func(res http.ResponseWriter, req *http.Request) {
		tenantExtRef := mux.Vars(req)["external_reference"]

		body := req.Body
		defer body.Close()

		var request account.Account
		if err := json.NewDecoder(body).Decode(&request); err != nil {
			log.WithError(err).Error("unable to read create account request")
			common.WriteUnreadableRequestError(res)
			return
		}
		created, err := accountSvc.Create(tenantExtRef, request)
		if err != nil {
			log.WithError(err).Error("unable to add new account")
			common.WriteFileServiceError(res, err)
			return
		}

		common.WriteOKWithEntity(res, created)

		res.WriteHeader(status)
		if err := json.NewEncoder(res).Encode(entity); err != nil {
			log.WithError(err).Error("unable to send response, the entity could not be encoded")
		}
	}
}
