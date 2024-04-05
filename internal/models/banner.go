package models

type Banner struct {
	Id        int     `json:"banner_id" db:"id"`
	Tags      []int32 `json:"tag_ids"`
	Feature   int32   `json:"feature_id" db:"feature_id"`
	Content   string  `json:"content" db:"content"`
	IsActive  bool    `json:"is_active" db:"is_active"`
	CreatedAt bool    `json:"created_at" db:"created_at"`
	UpdatedAt bool    `json:"updated_at" db:"updated_at"`
}

type BannerTags struct {
	Id        int    `json:"banner_id" db:"id"`
	Tag       int32  `json:"tag_ids"`
	Feature   int32  `json:"feature_id" db:"feature_id"`
	Content   string `json:"content" db:"content"`
	IsActive  bool   `json:"is_active" db:"is_active"`
	CreatedAt bool   `json:"created_at" db:"created_at"`
	UpdatedAt bool   `json:"updated_at" db:"updated_at"`
}
