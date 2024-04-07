package http

import (
	"avito_intern/internal/transport/http/response"
	"encoding/json"
	"go.uber.org/zap"
	"net/http"
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
