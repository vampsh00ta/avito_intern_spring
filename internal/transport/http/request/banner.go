package request

type GetBannerForUser struct {
	TagID           int32 `json:"tag_id" validate:"required" schema:"tag_id"`
	FeatureID       int32 `json:"feature_id" validate:"required" schema:"feature_id"`
	UseLastRevision bool  `json:"use_last_revision" schema:"use_last_revision"`
}

type GetBanners struct {
	TagID     int32 `json:"tag_id" schema:"tag_id"`
	FeatureID int32 `json:"feature_id" schema:"feature_id"`
	Limit     int32 `json:"limit" schema:"limit"`
	Offset    int32 `json:"offset" schema:"offset"`
}
