package core

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

//BadRequestError http status bad request error
type BadRequestError struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
}

type HttpErrorFmt interface {
	Error() string
}

func (e *BadRequestError) Error() string {
	return strconv.FormatInt(int64(e.Code), 10) + " : " + e.Message
}

type HttpStatusError struct {
	Status int
	Body   string
}

func (e HttpStatusError) Error() string {
	return fmt.Sprintf("%d:%s", e.Status, e.Body)
}

func HttpUnauthorizedError() error {
	return HttpStatusError{http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized)}
}

func HttpNotFoundError() error {
	return HttpStatusError{http.StatusNotFound, http.StatusText(http.StatusNotFound)}
}

func HttpOKError() error {
	return HttpStatusError{http.StatusOK, http.StatusText(http.StatusOK)}
}

func HttpNoContentError() error {
	return HttpStatusError{http.StatusNoContent, http.StatusText(http.StatusNoContent)}
}

func HttpInternalServerError() error {
	return HttpStatusError{http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError)}
}

func HttpBadRequestCustomError(code int32, message string) error {
	badR := BadRequestError{code, message}
	jsonString, _ := json.Marshal(&badR)
	return HttpStatusError{http.StatusBadRequest, string(jsonString)}
}

func HttpBadRequestError(httpError HttpErrorFmt) error {
	return HttpStatusError{http.StatusBadRequest, httpError.Error()}
}
