package http

import (
	"avito_intern/internal/transport/http/request"
	"avito_intern/internal/transport/http/response"
	"net/http"
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
