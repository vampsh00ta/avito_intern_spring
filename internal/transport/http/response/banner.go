package response

import "avito_intern/internal/models"

type GetBannerForUser struct {
	Content string `json:"content"`
}
type GetBanners []models.Banner
type GetBannerHistory []models.Banner

//	type GetBanners struct {
//		Banners []models.Banner `json:"banners"`
//	}
type CreateBanner struct {
	Id int `json:"id"  `
}
type DeleteBannerByTagAndFeature struct {
	Id int `json:"id"  `
}
