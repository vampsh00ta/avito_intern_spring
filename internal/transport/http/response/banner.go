package response

import "avito_intern/internal/models"

type ArrayBanners []models.Banner

type GetBannerForUser struct {
	Content string `json:"content"`
}
type GetBanners struct {
	ArrayBanners
}
