package http

import (
	"avito_intern/internal/errs"
	"avito_intern/internal/models"
	"avito_intern/internal/transport/http/request"
	"avito_intern/internal/transport/http/response"
	"encoding/json"
	"fmt"
	"net/http"
)

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
	if err != nil {
		t.handleHTTPError(w, err, methodName, http.StatusInternalServerError)

		return
	}
	if !userBanner.IsActive {
		userBanner.Content = ""
	}
	t.handleHTTPOk(w, response.GetBannerForUser{Content: userBanner.Content}, methodName, http.StatusOK)
}
func (t transport) GetBanners(w http.ResponseWriter, r *http.Request) {
	methodName := "GetBanners"

	if code, err := t.permission(w, r, models.Admin); err != nil {
		t.handleHTTPError(w, err, methodName, code)
		return
	}
	var req request.GetBanners

	if err := decoder.Decode(&req, r.URL.Query()); err != nil {
		t.handleHTTPError(w, err, methodName, http.StatusBadRequest)
		return
	}

	banners, err := t.s.GetBanners(r.Context(), req.TagID, req.FeatureID, req.Limit, req.Offset)
	if err != nil {
		t.handleHTTPError(w, err, methodName, http.StatusInternalServerError)

		return
	}
	t.handleHTTPOk(w, response.GetBanners(banners), methodName, http.StatusOK)
}

func (t transport) DeleteBannerByID(w http.ResponseWriter, r *http.Request) {
	methodName := "DeleteBannerByID"

	if code, err := t.permission(w, r, models.Admin); err != nil {
		t.handleHTTPError(w, err, methodName, code)
		return
	}

	ID, err := getIdFromUrl(r)
	if err != nil {
		t.handleHTTPError(w, err, methodName, http.StatusBadRequest)
		return
	}

	if err := t.s.DeleteBannerByID(r.Context(), ID); err != nil {
		t.handleHTTPError(w, err, methodName, http.StatusInternalServerError)

		return
	}
	t.handleHTTPOk(w, nil, methodName, http.StatusCreated)
}
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
		t.handleHTTPError(w, fmt.Errorf(errs.IncorrectJSONErr), methodName, http.StatusBadRequest)
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
	t.handleHTTPOk(w, response.CreateBanner{id}, methodName, http.StatusCreated)
}
func (t transport) ChangeBanner(w http.ResponseWriter, r *http.Request) {
	methodName := "ChangeBanner"

	if code, err := t.permission(w, r, models.Admin); err != nil {
		t.handleHTTPError(w, err, methodName, code)
		return
	}

	ID, err := getIdFromUrl(r)
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
			t.handleHTTPError(w, fmt.Errorf(errs.IncorrectJSONErr), methodName, http.StatusBadRequest)
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

func (t transport) DeleteBannerByTagAndFeature(w http.ResponseWriter, r *http.Request) {
	methodName := "DeleteBannerByTagAndFeature"

	if code, err := t.permission(w, r, models.Admin); err != nil {
		t.handleHTTPError(w, err, methodName, code)
		return
	}

	var req request.DeleteBannerByTagAndFeature

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
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
	t.handleHTTPOk(w, response.DeleteBannerByTagAndFeature{id}, methodName, http.StatusCreated)
}

func (t transport) GetBannerWithHistory(w http.ResponseWriter, r *http.Request) {
	methodName := "GetBannerWithHistory"

	if code, err := t.permission(w, r, models.Admin); err != nil {
		t.handleHTTPError(w, err, methodName, code)
		return
	}

	ID, err := getIdFromUrl(r)
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
	t.handleHTTPOk(w, banners, methodName, http.StatusOK)
}
