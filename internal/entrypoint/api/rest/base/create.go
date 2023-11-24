package base

import (
	"net/http"

	log "github.com/sirupsen/logrus"

	"gitlab.com/amihan/core/base.git/internal/component/user"
	"gitlab.com/amihan/core/base.git/internal/entrypoint/api/rest/common"
)

func createUser(user user.Service) http.HandlerFunc {
	log.Info("handler registered: user::create")

	return func(res http.ResponseWriter, req *http.Request) {
		common.WriteOK(res, "OK")
	}
}
