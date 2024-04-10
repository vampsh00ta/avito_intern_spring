package response

import "avito_intern/internal/models"

type GetBannerForUser struct {
	Content string `json:"content"`
}
type (
	GetBanners       []models.Banner
	GetBannerHistory []models.Banner
)

//	type GetBanners struct {
//		Banners []models.Banner `json:"banners"`
//	}
type CreateBanner struct {
	ID int `json:"id"  `
}
type DeleteBannerByTagAndFeature struct {
	ID int `json:"id" `
}
