package http

import (
	"avito_intern/internal/transport/http/request"
	"avito_intern/internal/transport/http/response"
	"encoding/json"
	"net/http"
)

func (t transport) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	methodName := "Login"
	var user request.Login

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		t.handleHTTPError(w, err, methodName, http.StatusBadRequest)

		return
	}
	if err := validate.Struct(user); err != nil {
		t.handleHTTPError(w, err, methodName, http.StatusBadRequest)

		return
	}
	jwtToken, err := t.s.Login(r.Context(), user.Username)
	if err != nil {
		t.handleHTTPError(w, err, methodName, http.StatusInternalServerError)

		return
	}
	t.handleHTTPOk(w, response.Login{Access: jwtToken}, methodName, http.StatusOK)
}
