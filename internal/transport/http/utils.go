package http

import (
	"avito_intern/internal/errs"
	"avito_intern/internal/transport/http/response"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"go.uber.org/zap"
)

func (t transport) handleHTTPError(w http.ResponseWriter, err error, method string, status int) {
	w.Header().Set("Content-Type", "application/json")

	logError := err
	err = errs.Handle(err)
	if errors.Is(err, errs.NoRowsInResultErr) {
		fmt.Println(err)
		status = http.StatusNotFound
	}
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response.Error{Error: err.Error()})
	t.l.Error(method, zap.Error(logError))
}

func (t transport) handleHTTPOk(w http.ResponseWriter, resp interface{}, method string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if resp != nil {
		json.NewEncoder(w).Encode(resp)
	}
	t.l.Info(method)
}

func IsJSON(str string) bool {
	var js json.RawMessage
	return json.Unmarshal([]byte(str), &js) == nil
}

func getIDFromURL(r *http.Request) (int, error) {
	strID := r.PathValue("id")
	if strID == "" {
		return -1, errs.NilIDErr
	}
	ID, err := strconv.Atoi(strID)
	if err != nil {
		return -1, errs.WrongIDErr
	}
	return ID, nil
}

//	func (t transport) permission(w http.ResponseWriter, r *http.Request, groupIDs ...int) error {
//		jwtToken := r.Header.Get("Authorization")
//		admin, errs := t.s.IsLogged(r.Context(), jwtToken)
//		if errs != nil {
//			return errs
//		}
//		if !admin {
//			return fmt.Errorf(models.NotAdminErr)
//		}
//		return nil
//	}
//
//	func (t transport) adminPermission(w http.ResponseWriter, r *http.Request) error {
//		jwtToken := r.Header.Get("Authorization")
//		admin, errs := t.s.IsAdmin(r.Context(), jwtToken)
//		if errs != nil {
//			return errs
//		}
//		if !admin {
//			return fmt.Errorf(models.NotAdminErr)
//		}
//		return nil
//	}
func (t transport) permission(_ http.ResponseWriter, r *http.Request, groupIDs ...int) (int, error) {
	jwtToken := r.Header.Get("Authorization")
	ok, err := t.s.Permission(r.Context(), jwtToken, groupIDs...)
	if err != nil {
		return http.StatusUnauthorized, err
	}
	if !ok {
		return http.StatusForbidden, errs.WrongRoleErr
	}

	return http.StatusAccepted, nil
}
