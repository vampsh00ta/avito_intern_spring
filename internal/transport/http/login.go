package http

import (
	"avito_intern/internal/transport/http/request"
	"avito_intern/internal/transport/http/response"
	"encoding/json"
	"net/http"
)

// @Summary     Login
// @Description Авторизация
// @Tags        Login
// @Accept      json
// @Produce     json
// @Param data body request.Login true "Модель запроса"
// @Success     201 {object} response.Login "Access token"
// @Failure     400 {object} response.Error Некорректные данные
// @Failure     401 {object} response.Error Пользователь не авторизован
// @Failure     403 {object} response.Error Пользователь не имеет доступа
// @Failure     404 {object} response.Error Пользователь не найден
// @Failure     500 {object} response.Error Внутренняя ошибка сервера
// @Security ApiKeyAuth
// @Router      /login [post]
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
