package common

import "net/http"

// WriteOK writes an ok response with the provided entity
func WriteOK(res http.ResponseWriter, entity interface{}) {
	WriteResponse(res, http.StatusOK, entity)
}

// WriteCreated writes a CREATED response with no entity
func WriteCreated(res http.ResponseWriter) {
	res.WriteHeader(http.StatusCreated)
}

// WriteRepositoryGatewayTimeout is when the repository cannot be accessed
func WriteRepositoryGatewayTimeout(res http.ResponseWriter) {
	status := http.StatusGatewayTimeout
	err := ErrorBuilder().
		Status(status).
		Title("Repository Inaccessible").
		Detail("Repository is either inaccessible or down").
		Build()
	WriteErrorResponse(res, status, err)
}

// WriteUnknownError is when you don't know what the error is, this shouldn't be the case and is a bug in the code
func WriteUnknownError(res http.ResponseWriter) {
	status := http.StatusInternalServerError
	err := ErrorBuilder().
		Status(status).
		Title("Unknown Error").
		Detail("An unknown error has occurred, contact the system administrator").
		Build()
	WriteErrorResponse(res, status, err)
}

// WriteUnreadableRequestError is when the request body cannot be pased into JSON
func WriteUnreadableRequestError(res http.ResponseWriter) {
	status := http.StatusBadRequest
	jsonErr := ErrorBuilder().
		Status(status).
		Title("Request cannot be read").
		Detail("Request is not a valid JSON format").
		SourcePointer("/").
		Build()
	WriteErrorResponse(res, status, jsonErr)
}
