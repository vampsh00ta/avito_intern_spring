package http

import (
	"avito_intern/internal/models"
	"avito_intern/internal/transport/http/request"
	"avito_intern/internal/transport/http/response"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func (t transport) GetBannerForUser(w http.ResponseWriter, r *http.Request) {
	methodName := "GetBannerForUser"

	var req request.GetBannerForUser

	if err := decoder.Decode(&req, r.URL.Query()); err != nil {
		t.handleError(w, err, err, methodName, http.StatusBadRequest)
		return
	}
	if err := validate.Struct(req); err != nil {
		t.handleError(w, err, err, methodName, http.StatusBadRequest)
		return
	}

	userBanner, err := t.s.GetBannerForUser(r.Context(), req.TagID, req.FeatureID, req.UseLastRevision)
	if err != nil {
		t.handleError(w, err, err, methodName, http.StatusInternalServerError)

		return
	}
	t.handleOk(w, response.GetBannerForUser{Content: userBanner.Content}, methodName, http.StatusOK)
}
func (t transport) GetBanners(w http.ResponseWriter, r *http.Request) {
	methodName := "GetBanners"

	var req request.GetBanners

	if err := decoder.Decode(&req, r.URL.Query()); err != nil {
		t.handleError(w, err, err, methodName, http.StatusBadRequest)
		return
	}

	banners, err := t.s.GetBanners(r.Context(), req.TagID, req.FeatureID, req.Limit, req.Offset)
	if err != nil {
		t.handleError(w, err, err, methodName, http.StatusInternalServerError)

		return
	}
	t.handleOk(w, banners, methodName, http.StatusOK)
}

func (t transport) DeleteBannerByID(w http.ResponseWriter, r *http.Request) {
	methodName := "DeleteBannerByID"
	strID := r.PathValue("id")
	if strID == "" {
		t.handleError(w, fmt.Errorf("nil id"), fmt.Errorf("nil id"), methodName, http.StatusBadRequest)
		return
	}
	ID, err := strconv.Atoi(strID)
	if err != nil {
		t.handleError(w, fmt.Errorf("wrong id"), fmt.Errorf("wrong id"), methodName, http.StatusBadRequest)
		return
	}

	if err := t.s.DeleteBannerByID(r.Context(), ID); err != nil {
		t.handleError(w, err, err, methodName, http.StatusInternalServerError)

		return
	}
	t.handleOk(w, nil, methodName, http.StatusCreated)
}
func (t transport) CreateBanner(w http.ResponseWriter, r *http.Request) {
	methodName := "CreateBanner"

	var req request.CreateBanner

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		t.handleError(w, err, err, methodName, http.StatusBadRequest)

		return
	}
	if err := validate.Struct(req); err != nil {
		t.handleError(w, err, err, methodName, http.StatusBadRequest)
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
		t.handleError(w, err, err, methodName, http.StatusInternalServerError)

		return
	}
	t.handleOk(w, response.CreateBanner{id}, methodName, http.StatusCreated)
}
func (t transport) ChangeBanner(w http.ResponseWriter, r *http.Request) {
	methodName := "ChangeBanner"
	ID, err := getIdFromUrl(r)
	if err != nil {
		t.handleError(w, err, err, methodName, http.StatusBadRequest)
		return
	}
	var req request.ChangeBanner
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		t.handleError(w, err, err, methodName, http.StatusBadRequest)

		return
	}

	banner := models.BannerChange{
		Tags:     req.Tags,
		Feature:  req.Feature,
		Content:  req.Content,
		IsActive: req.IsActive,
	}
	fmt.Println(r.Context())
	if err := t.s.ChangeBanner(r.Context(), ID, banner); err != nil {
		t.handleError(w, err, err, methodName, http.StatusInternalServerError)

		return
	}
	t.handleOk(w, nil, methodName, http.StatusCreated)
}
