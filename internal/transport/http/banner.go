package http

import (
	"avito_intern/internal/errs"
	"avito_intern/internal/models"
	"avito_intern/internal/transport/http/request"
	"avito_intern/internal/transport/http/response"
	"encoding/json"
	"net/http"
)

// @Summary     GetBannerForUser
// @Description Получение баннера для пользователя
// @Tags        Banner
// @Accept      json
// @Param tag_id query string true "Тэг пользователя"
// @Param feature_id query string true "Идентификатор фичи"
// @Param use_last_revision query string false "Получать актуальную информацию "
// @Produce     json
// @Success     200 {object} response.GetBannerForUser
// @Failure     400 {object} response.Error Некорректные данные
// @Failure     401 {object} response.Error Пользователь не авторизован
// @Failure     403 {object} response.Error Пользователь не имеет доступа
// @Failure     404 {object} response.Error Баннер не найден
// @Failure     500 {object} response.Error Внутренняя ошибка сервера
// @Security ApiKeyAuth
// @Router      /user_banner [get].
func (t transport) GetBannerForUser(w http.ResponseWriter, r *http.Request) {
	methodName := "GetBannerForUser"

	if code, err := t.permission(w, r, models.OrdinaryUser, models.Admin); err != nil {
		t.handleHTTPError(w, err, methodName, code)
		return
	}
	var req request.GetBannerForUser

	if err := decoder.Decode(&req, r.URL.Query()); err != nil {
		t.handleHTTPError(w, err, methodName, http.StatusBadRequest)

		return
	}

	if err := validate.Struct(req); err != nil {
		t.handleHTTPError(w, err, methodName, http.StatusBadRequest)

		return
	}
	userBanner, err := t.s.GetBannerForUser(r.Context(), req.TagID, req.FeatureID, req.UseLastRevision)
	if !userBanner.IsActive {
		userBanner.Content = ""
	}
	if err != nil {
		t.handleHTTPError(w, err, methodName, http.StatusInternalServerError)
		return
	}
	if !userBanner.IsActive {
		t.handleHTTPError(w, errs.NoRowsInResult, methodName, http.StatusInternalServerError)
		return
	}

	t.handleHTTPOk(w, response.GetBannerForUser{Content: userBanner.Content}, methodName, http.StatusOK)
}

// @Summary     GetBanners
// @Description Получение всех баннеров c фильтрацией по фиче и/или тегу
// @Tags        Banner
// @Accept      json
// @Param feature_id query string false "Идентификатор фичи"
// @Param tag_id query string false "Идентификатор тега"
// @Param limit query string false "Лимит"
// @Param offset query string false "Оффсет"
// @Produce     json
// @Success     200 {object} response.GetBanners
// @Failure     400 {object} response.Error Некорректные данные
// @Failure     401 {object} response.Error Пользователь не авторизован
// @Failure     403 {object} response.Error Пользователь не имеет доступа
// @Failure     404 {object} response.Error Баннеры не найдены
// @Failure     500 {object} response.Error Внутренняя ошибка сервера
// @Security ApiKeyAuth
// @Router      /banner [get].
func (t transport) GetBanners(w http.ResponseWriter, r *http.Request) {
	methodName := "GetBanners"

	if code, err := t.permission(w, r, models.Admin); err != nil {

		t.handleHTTPError(w, err, methodName, code)
		return
	}
	var req request.GetBanners

	if err := decoder.Decode(&req, r.URL.Query()); err != nil {
		t.handleHTTPError(w, errs.Validation, methodName, http.StatusBadRequest)

		return
	}

	banners, err := t.s.GetBanners(r.Context(), req.TagID, req.FeatureID, req.Limit, req.Offset)
	if err != nil {
		t.handleHTTPError(w, err, methodName, http.StatusInternalServerError)

		return
	}
	t.handleHTTPOk(w, response.GetBanners(banners), methodName, http.StatusOK)
}

// @Summary     DeleteBannerByID
// @Description Удаление баннера по идентификатору
// @Tags        Banner
// @Produce     json
// @Param        id   path      int  true  "Идентификатор баннера"
// @Success     204  "Баннер успешно удален"
// @Failure     400  {object} response.Error "Некорректные данные"
// @Failure     401  {object} response.Error "Пользователь не авторизован"
// @Failure     403  {object} response.Error "Пользователь не имеет доступа"
// @Failure     404  {object} response.Error "Баннер не найден"
// @Failure     500  {object} response.Error "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
// @Router      /banner/{id} [delete].
func (t transport) DeleteBannerByID(w http.ResponseWriter, r *http.Request) {
	methodName := "DeleteBannerByID"

	if code, err := t.permission(w, r, models.Admin); err != nil {
		t.handleHTTPError(w, err, methodName, code)
		return
	}

	ID, err := getIDFromURL(r)
	if err != nil {
		t.handleHTTPError(w, err, methodName, http.StatusBadRequest)
		return
	}

	if err := t.s.DeleteBannerByID(r.Context(), ID); err != nil {
		t.handleHTTPError(w, err, methodName, http.StatusInternalServerError)

		return
	}
	t.handleHTTPOk(w, nil, methodName, http.StatusNoContent)
}

// @Summary     CreateBanner
// @Description Создание нового баннера
// @Tags        Banner
// @Accept      json
// @Produce     json
// @Param data body request.CreateBanner true "Модель запроса"
// @Success     201 {object} response.CreateBanner "Created"
// @Failure     400 {object} response.Error Некорректные данные
// @Failure     401 {object} response.Error Пользователь не авторизован
// @Failure     404 {object} response.Error Баннер для тэга не найден
// @Failure     500 {object} response.Error Внутренняя ошибка сервера
// @Security ApiKeyAuth
// @Router      /banner [post].
func (t transport) CreateBanner(w http.ResponseWriter, r *http.Request) {
	methodName := "CreateBanner"

	if code, err := t.permission(w, r, models.Admin); err != nil {
		t.handleHTTPError(w, err, methodName, code)
		return
	}

	var req request.CreateBanner

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		t.handleHTTPError(w, err, methodName, http.StatusBadRequest)

		return
	}
	if err := validate.Struct(req); err != nil {
		t.handleHTTPError(w, err, methodName, http.StatusBadRequest)
		return
	}
	if ok := json.Valid([]byte(req.Content)); !ok {
		t.handleHTTPError(w, errs.IncorrectJSON, methodName, http.StatusBadRequest)
		return
	}
	banner := models.Banner{
		Tags:     req.Tags,
		Feature:  req.Feature,
		Content:  req.Content,
		IsActive: req.IsActive,
	}
	id, err := t.s.CreateBanner(r.Context(), banner)
	if err != nil {
		t.handleHTTPError(w, err, methodName, http.StatusInternalServerError)

		return
	}
	t.handleHTTPOk(w, response.CreateBanner{ID: id}, methodName, http.StatusCreated)
}

// @Summary     ChangeBanner
// @Description Обновление содержимого баннера
// @Tags        Banner
// @Accept      json
// @Produce     json
// @Param        id   path      int  true  "Идентификатор баннера"
// @Param data body request.ChangeBanner true "Модель запроса"
// @Success     201 {object} response.CreateBanner "Created"
// @Failure     400 {object} response.Error Некорректные данные
// @Failure     401 {object} response.Error Пользователь не авторизован
// @Failure     403 {object} response.Error Пользователь не имеет доступа
// @Failure     404 {object} response.Error Баннер не найден
// @Failure     500 {object} response.Error Внутренняя ошибка сервера
// @Security ApiKeyAuth
// @Router      /banner/{id} [patch].
func (t transport) ChangeBanner(w http.ResponseWriter, r *http.Request) {
	methodName := "ChangeBanner"

	if code, err := t.permission(w, r, models.Admin); err != nil {
		t.handleHTTPError(w, err, methodName, code)
		return
	}

	ID, err := getIDFromURL(r)
	if err != nil {
		t.handleHTTPError(w, err, methodName, http.StatusBadRequest)
		return
	}
	var req request.ChangeBanner
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		t.handleHTTPError(w, err, methodName, http.StatusBadRequest)

		return
	}

	if req.Content != nil {
		if ok := IsJSON(*req.Content); !ok {
			t.handleHTTPError(w, errs.IncorrectJSON, methodName, http.StatusBadRequest)
			return
		}
	}

	banner := models.BannerChange{
		Tags:     req.Tags,
		Feature:  req.Feature,
		Content:  req.Content,
		IsActive: req.IsActive,
	}
	if err := t.s.ChangeBanner(r.Context(), ID, banner); err != nil {
		t.handleHTTPError(w, err, methodName, http.StatusInternalServerError)

		return
	}
	t.handleHTTPOk(w, nil, methodName, http.StatusCreated)
}

// @Summary     DeleteBannerByTagAndFeature
// @Description Удаление банера по  тэгу и фиче
// @Tags        Banner
// @Produce     json
// @Param tag_id query int true "Идентификатор тега"
// @Param feature_id query int true "Идентификатор фичи"
// @Success     204 {object} response.DeleteBannerByTagAndFeature "Deleted"
// @Failure     400 {object} response.Error "Некорректные данные"
// @Failure     401 {object} response.Error "Пользователь не авторизован"
// @Failure     403 {object} response.Error "Пользователь не имеет доступа"
// @Failure     404 {object} response.Error "Баннер не найден"
// @Failure     500 {object} response.Error "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
// @Router      /banner [delete].
func (t transport) DeleteBannerByTagAndFeature(w http.ResponseWriter, r *http.Request) {
	methodName := "DeleteBannerByTagAndFeature"

	if code, err := t.permission(w, r, models.Admin); err != nil {
		t.handleHTTPError(w, err, methodName, code)
		return
	}

	var req request.DeleteBannerByTagAndFeature
	if err := decoder.Decode(&req, r.URL.Query()); err != nil {
		t.handleHTTPError(w, err, methodName, http.StatusBadRequest)

		return
	}

	if err := validate.Struct(req); err != nil {
		t.handleHTTPError(w, err, methodName, http.StatusBadRequest)
		return
	}

	id, err := t.s.DeleteBannerByTagAndFeature(r.Context(), req.FeatureID, req.TagID)
	if err != nil {
		t.handleHTTPError(w, err, methodName, http.StatusInternalServerError)
		return
	}
	t.handleHTTPOk(w, response.DeleteBannerByTagAndFeature{ID: id}, methodName, http.StatusNoContent)
}

// @Summary     GetBannerWithHistory
// @Description История изменений банера
// @Tags        Banner
// @Accept      json
// @Produce     json
// @Param        id   path      int  true  "Идентификатор баннера"
// @Param 		limit query int false "Лимит (до 3 )"
// @Success     200 {object} response.GetBannerHistory "История"
// @Failure     400 {object} response.Error Некорректные данные
// @Failure     401 {object} response.Error Пользователь не авторизован
// @Failure     403 {object} response.Error Пользователь не имеет доступа
// @Failure     404 {object} response.Error История  не найдена
// @Failure     500 {object} response.Error Внутренняя ошибка сервера
// @Security ApiKeyAuth
// @Router      /banner_history/{id} [get].
func (t transport) GetBannerWithHistory(w http.ResponseWriter, r *http.Request) {
	methodName := "GetBannerWithHistory"

	if code, err := t.permission(w, r, models.Admin); err != nil {
		t.handleHTTPError(w, err, methodName, code)
		return
	}

	ID, err := getIDFromURL(r)
	if err != nil {
		t.handleHTTPError(w, err, methodName, http.StatusBadRequest)
		return
	}
	var req request.GetBannerHistory
	if err := decoder.Decode(&req, r.URL.Query()); err != nil {
		t.handleHTTPError(w, err, methodName, http.StatusBadRequest)

		return
	}

	banners, err := t.s.GetBannerWithHistory(r.Context(), ID, req.Limit)
	if err != nil {
		t.handleHTTPError(w, err, methodName, http.StatusInternalServerError)

		return
	}

	t.handleHTTPOk(w, response.GetBannerHistory(banners), methodName, http.StatusOK)
}
