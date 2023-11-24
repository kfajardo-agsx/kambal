package common

import (
	"net/http"

	"github.com/kfajardo-agsx/kambal.git/internal/component/common"
)

func WriteFileServiceError(res http.ResponseWriter, err error) {
	var status int
	var errResponse JSONError

	switch errNew := err.(type) {
	case *common.APIError:
		switch errNew.Type {
		case common.ErrorTypeUnreachable:
			status = http.StatusBadGateway
			errResponse = ErrorBuilder().
				Status(status).
				Title("Unreachable").
				Detail(errNew.Error()).
				Build()
		case common.ErrorTypeNotFound:
			status = http.StatusNotFound
			errResponse = ErrorBuilder().
				Status(status).
				Title("Not found").
				Detail(errNew.Error()).
				Build()
		case common.ErrorTypeConflict:
			status = http.StatusConflict
			errResponse = ErrorBuilder().
				Status(status).
				Title("Conflict").
				Detail(errNew.Error()).
				Build()
		case common.ErrorTypeBadRequest:
			status = http.StatusBadRequest
			errResponse = ErrorBuilder().
				Status(status).
				Title("Bad Request").
				Detail(errNew.Error()).
				Build()
		case common.ErrorTypeUnauthorized:
			status = http.StatusUnauthorized
			errResponse = ErrorBuilder().
				Status(status).
				Title("Unauthorized").
				Detail(errNew.Error()).
				Build()
		case common.ErrorTypeForbidden:
			status = http.StatusForbidden
			errResponse = ErrorBuilder().
				Status(status).
				Title("Forbidden").
				Detail(errNew.Error()).
				Build()
		case common.ErrorTypeServiceUnavailable:
			status = http.StatusServiceUnavailable
			errResponse = ErrorBuilder().
				Status(status).
				Title("Service unavailable").
				Detail(errNew.Error()).
				Build()
		case common.ErrorExternalService:
			status = http.StatusInternalServerError
			errResponse = ErrorBuilder().
				Status(status).
				Title("External Service Error").
				Detail(errNew.Error()).
				Build()
		default:
			WriteUnknownError(res)
			return
		}
		WriteErrorResponse(res, status, errResponse)
	default:
		WriteUnknownError(res)
	}
}
