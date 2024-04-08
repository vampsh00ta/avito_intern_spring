package http

import (
	"avito_intern/internal/models"
	"avito_intern/internal/transport/http/response"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

func handleError(err error) error {

	if pgErr, ok := err.(*pgconn.PgError); ok {
		switch pgErr.Code {
		case "23505":
			err = fmt.Errorf(models.DublicateErr)
		}
	}
	return err
}
func (t transport) handleHTTPError(w http.ResponseWriter, err error, method string, status int) {
	w.WriteHeader(status)
	logError := err
	err = handleError(err)
	json.NewEncoder(w).Encode(response.Error{Error: err.Error()})
	t.l.Error(method, zap.Error(logError))
}
func (t transport) handleHTTPOk(w http.ResponseWriter, resp interface{}, method string, status int) {

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
func getIdFromUrl(r *http.Request) (int, error) {
	strID := r.PathValue("id")
	if strID == "" {
		return -1, fmt.Errorf("nil id")
	}
	ID, err := strconv.Atoi(strID)
	if err != nil {
		return -1, fmt.Errorf("wrong id")
	}
	return ID, nil
}

//	func (t transport) permission(w http.ResponseWriter, r *http.Request, groupIDs ...int) error {
//		jwtToken := r.Header.Get("Authorization")
//		admin, err := t.s.IsLogged(r.Context(), jwtToken)
//		if err != nil {
//			return err
//		}
//		if !admin {
//			return fmt.Errorf(models.NotAdminErr)
//		}
//		return nil
//	}
//
//	func (t transport) adminPermission(w http.ResponseWriter, r *http.Request) error {
//		jwtToken := r.Header.Get("Authorization")
//		admin, err := t.s.IsAdmin(r.Context(), jwtToken)
//		if err != nil {
//			return err
//		}
//		if !admin {
//			return fmt.Errorf(models.NotAdminErr)
//		}
//		return nil
//	}
func (t transport) permission(w http.ResponseWriter, r *http.Request, groupIDs ...int) (int, error) {
	jwtToken := r.Header.Get("Authorization")
	ok, err := t.s.Permission(r.Context(), jwtToken, groupIDs...)
	if err != nil {
		return http.StatusUnauthorized, err
	}
	if !ok {
		return http.StatusForbidden, fmt.Errorf(models.WrongRoleErr)

	}

	return http.StatusAccepted, nil
}
