package response

import "avito_intern/internal/models"

type GetBannerForUser struct {
	Content string `json:"content"`
}
type GetBanners struct {
	Banners []models.Banner `json:"banners"`
}
type CreateBanner struct {
	Id int `json:"id"  `
}
