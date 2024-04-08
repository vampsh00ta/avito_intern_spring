package http

import (
	"avito_intern/internal/transport/http/response"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

func (t transport) handleError(w http.ResponseWriter, err, handledError error, method string, status int) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response.Error{Error: handledError.Error()})
	t.l.Error(method, zap.Error(err))
}
func (t transport) handleOk(w http.ResponseWriter, resp interface{}, method string, status int) {

	w.WriteHeader(status)
	json.NewEncoder(w).Encode(resp)
	t.l.Info(method)
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
