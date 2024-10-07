package router

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ashihara-api/core/domain/errors"
	"github.com/ashihara-api/core/interface/transport/http/render"
)

// BindFromJSON ...
func BindFromJSON(body io.ReadCloser, v interface{}) (err error) {
	if err := json.NewDecoder(body).Decode(v); err != nil {
		return errors.NewCause(err, errors.CaseBadRequest)
	}
	return nil
}

// RenderJSON ...
func RenderJSON(w http.ResponseWriter, s int, v interface{}) {
	if err := render.JSON(w, s, v); err != nil {
		err := errors.NewCause(fmt.Errorf("fail to encode response : %w", err), errors.CaseBackendError)
		render.JSON(w, http.StatusInternalServerError, err) // nolint: errcheck
		return
	}
}

// RenderError ...
func RenderError(w http.ResponseWriter, err error) {
	var ec *errors.Cause
	if !errors.As(err, &ec) {
		errors.As(errors.NewCause(err, errors.CaseBackendError), &ec)
	}
	render.JSON(w, ec.Code(), ec) // nolint: errcheck
}
